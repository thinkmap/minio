/*
 * MinIO Cloud Storage, (C) 2018 MinIO, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package madmin

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

// QuotaType represents bucket quota type
type QuotaType string

const (
	// HardQuota specifies a hard quota of usage for bucket
	HardQuota QuotaType = "hard"
	// FIFOQuota specifies a quota limit beyond which older files are deleted from bucket
	FIFOQuota QuotaType = "fifo"
)

// IsValid returns true if quota type is one of FIFO or Hard
func (t QuotaType) IsValid() bool {
	return t == HardQuota || t == FIFOQuota
}

// BucketQuota holds bucket quota restrictions
type BucketQuota struct {
	Quota uint64    `json:"quota"`
	Type  QuotaType `json:"quotatype"`
}

// RemoveBucketQuota - removes quota config on a bucket.
func (adm *AdminClient) RemoveBucketQuota(ctx context.Context, bucket string) error {

	queryValues := url.Values{}
	queryValues.Set("bucket", bucket)

	reqData := requestData{
		relPath:     adminAPIPrefix + "/remove-bucket-quota",
		queryValues: queryValues,
	}

	// Execute DELETE on /minio/admin/v3/remove-bucket-quota to delete bucket quota.
	resp, err := adm.executeMethod(ctx, http.MethodDelete, reqData)
	defer closeResponse(resp)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return httpRespToErrorResponse(resp)
	}

	return nil
}

// GetBucketQuota - get info on a user
func (adm *AdminClient) GetBucketQuota(ctx context.Context, bucket string) (q BucketQuota, err error) {
	queryValues := url.Values{}
	queryValues.Set("bucket", bucket)

	reqData := requestData{
		relPath:     adminAPIPrefix + "/get-bucket-quota",
		queryValues: queryValues,
	}

	// Execute GET on /minio/admin/v3/get-quota
	resp, err := adm.executeMethod(ctx, http.MethodGet, reqData)

	defer closeResponse(resp)
	if err != nil {
		return q, err
	}

	if resp.StatusCode != http.StatusOK {
		return q, httpRespToErrorResponse(resp)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return q, err
	}
	if err = json.Unmarshal(b, &q); err != nil {
		return q, err
	}

	return q, nil
}

// SetBucketQuota - sets a bucket's quota.
func (adm *AdminClient) SetBucketQuota(ctx context.Context, bucket string, quota uint64, quotaType QuotaType) error {

	data, err := json.Marshal(BucketQuota{
		Quota: quota,
		Type:  quotaType,
	})
	if err != nil {
		return err
	}

	queryValues := url.Values{}
	queryValues.Set("bucket", bucket)

	reqData := requestData{
		relPath:     adminAPIPrefix + "/set-bucket-quota",
		queryValues: queryValues,
		content:     data,
	}

	// Execute PUT on /minio/admin/v3/set-bucket-quota to set quota for a bucket.
	resp, err := adm.executeMethod(ctx, http.MethodPut, reqData)

	defer closeResponse(resp)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return httpRespToErrorResponse(resp)
	}

	return nil
}
