package client

/*

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

import (
	"encoding/json"
	"strconv"

	"github.com/apache/trafficcontrol/pkg/tc"
)

// SetUserDeliveryService associates the given delivery services with the given user.
func (to *Session) SetDeliveryServiceUser(userID int, dses []int, replace bool) (*tc.UserDeliveryServicePostResponse, error) {
	uri := apiBase + `/deliveryservice_user`
	ds := tc.DeliveryServiceUserPost{UserID: &userID, DeliveryServices: &dses, Replace: &replace}
	jsonReq, err := json.Marshal(ds)
	if err != nil {
		return nil, err
	}
	resp := tc.UserDeliveryServicePostResponse{}
	_, err = post(to, uri, jsonReq, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteDeliveryServiceUser deletes the association between the given delivery service and user
func (to *Session) DeleteDeliveryServiceUser(userID int, dsID int) (*tc.UserDeliveryServiceDeleteResponse, error) {
	uri := apiBase + `/deliveryservice_user/` + strconv.Itoa(dsID) + `/` + strconv.Itoa(userID)
	resp := tc.UserDeliveryServiceDeleteResponse{}
	if _, err := del(to, uri, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
