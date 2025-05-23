/*	Copyright [2024] [Elia Renzoni]
*
*   Licensed under the Apache License, Version 2.0 (the "License");
*   you may not use this file except in compliance with the License.
*   You may obtain a copy of the License at
*
*      http://www.apache.org/licenses/LICENSE-2.0
*
*   Unless required by applicable law or agreed to in writing, software
*   distributed under the License is distributed on an "AS IS" BASIS,
*   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
*   See the License for the specific language governing permissions and
*   limitations under the License.
*/


/**
*	This file contains the implementation of an hash translator.
*	Given a hash generated by spooky_hash32.go translates it into the address of a bucket
*   bucket_id = hash % buffer_pool_size
*
**/
package store

const BUFFER_POOL_SIZE uint32 = 3000

type AddressBinder struct {
	bucketAddress uint32
}

func NewAddressBinder() *AddressBinder {
	return &AddressBinder{}
}

/**
*	@brief calculate the bucket address
*   @param hash value calculated by the hash function
*   @return bucket address
*   if the hash is for example 1177965712 the corresponding bucket
*   address is 712
 */
func (a *AddressBinder) TranslateHash(hash uint32) uint32 {
	a.bucketAddress = hash % BUFFER_POOL_SIZE
	return a.bucketAddress
}
