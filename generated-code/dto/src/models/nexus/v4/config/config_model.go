/*
 * Generated file models/nexus/v4/config/config_model.go.
 *
 * Product version: 1.0.0-SNAPSHOT
 *
 * Part of the GoLang Mock API - REST API for Mock Item Service
 *
 * (c) 2025 Nutanix Inc.  All rights reserved
 *
 */

/*
  Module nexus.v4.config of GoLang Mock API - REST API for Mock Item Service
*/
package config
import (
  import2 "github.com/nutanix/ntnx-api-golang-nexus-pc/generated-code/dto/models/common/v1/response"
  "encoding/json"
  "errors"
  "fmt"
  import1 "github.com/nutanix/ntnx-api-golang-nexus-pc/generated-code/dto/models/nexus/v4/error"
)
/*
Country information
*/
type Country struct {
  
  ObjectType_ *string `json:"$objectType,omitempty"`
  
  Reserved_ map[string]interface{} `json:"$reserved,omitempty"`
  
  UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
  /*
  State or province name
  */
  State *string `json:"state,omitempty"`
}

func (p *Country) MarshalJSON() ([]byte, error) {
  // Create Alias to avoid infinite recursion
  type Alias Country

  // Step 1: Marshal the known fields
  known, err := json.Marshal(Alias(*p))
  if err != nil {
  	return nil, err
  }

    // Step 2: Convert known to map for merging
    var knownMap map[string]interface{}
    if err := json.Unmarshal(known, &knownMap); err != nil {
    	return nil, err
    }
    delete(knownMap, "$unknownFields")
  
    // Step 3: Merge unknown fields
    for k, v := range p.UnknownFields_ {
    	knownMap[k] = v
    }
  
    // Step 4: Marshal final merged map
    return json.Marshal(knownMap)
}

func (p *Country) UnmarshalJSON(b []byte) error {
    // Step 1: Unmarshal into a generic map to capture all fields
    var allFields map[string]interface{}
	if err := json.Unmarshal(b, &allFields); err != nil {
		return err
	}

    // Step 2: Unmarshal into a temporary struct with known fields
	type Alias Country
	known := &Alias{}
	if err := json.Unmarshal(b, known); err != nil {
		return err
	}

    // Step 3: Assign known fields
	*p = *NewCountry()

    if known.ObjectType_ != nil {
        p.ObjectType_ = known.ObjectType_
    }
    if known.Reserved_ != nil {
        p.Reserved_ = known.Reserved_
    }
    if known.UnknownFields_ != nil {
        p.UnknownFields_ = known.UnknownFields_
    }
    if known.State != nil {
        p.State = known.State
    }

    // Step 4: Remove known JSON fields from allFields map
	delete(allFields, "$objectType")
	delete(allFields, "$reserved")
	delete(allFields, "$unknownFields")
	delete(allFields, "state")

    // Step 5: Assign remaining fields to UnknownFields_
	for key, value := range allFields {
      p.UnknownFields_[key] = value
    }

	return nil
}

func NewCountry() *Country {
  p := new(Country)
  p.ObjectType_ = new(string)
  *p.ObjectType_ = "nexus.v4.config.Country"
  p.Reserved_ = map[string]interface{}{"$fv": "v4.r1"}
  p.UnknownFields_ = map[string]interface{}{}



  return p
}



/*
Item entity for mock REST API
*/
type Item struct {
  
  ObjectType_ *string `json:"$objectType,omitempty"`
  
  Reserved_ map[string]interface{} `json:"$reserved,omitempty"`
  
  UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
  /*
  Description of the item
  */
  Description *string `json:"description,omitempty"`
  /*
  Unique identifier for the item
  */
  ItemId *int `json:"itemId,omitempty"`
  /*
  Path to item image file
  */
  ItemImageFile *string `json:"itemImageFile,omitempty"`
  /*
  Name of the item
  */
  ItemName *string `json:"itemName"`
  /*
  Type of item
  */
  ItemType *string `json:"itemType"`
  
  Location *Location `json:"location,omitempty"`
}

func (p *Item) MarshalJSON() ([]byte, error) {
  type ItemProxy Item

  // Step 1: Marshal known fields via proxy to enforce required fields
  baseStruct := struct {
    *ItemProxy
    ItemName *string `json:"itemName,omitempty"`
    ItemType *string `json:"itemType,omitempty"`
  }{
    ItemProxy : (*ItemProxy)(p),
    ItemName : p.ItemName,
    ItemType : p.ItemType,
  }

  known, err := json.Marshal(baseStruct)
  if err != nil {
  	return nil, err
  }

    // Step 2: Convert known to map for merging
    var knownMap map[string]interface{}
    if err := json.Unmarshal(known, &knownMap); err != nil {
    	return nil, err
    }
    delete(knownMap, "$unknownFields")
  
    // Step 3: Merge unknown fields
    for k, v := range p.UnknownFields_ {
    	knownMap[k] = v
    }
  
    // Step 4: Marshal final merged map
    return json.Marshal(knownMap)
}

func (p *Item) UnmarshalJSON(b []byte) error {
    // Step 1: Unmarshal into a generic map to capture all fields
    var allFields map[string]interface{}
	if err := json.Unmarshal(b, &allFields); err != nil {
		return err
	}

    // Step 2: Unmarshal into a temporary struct with known fields
	type Alias Item
	known := &Alias{}
	if err := json.Unmarshal(b, known); err != nil {
		return err
	}

    // Step 3: Assign known fields
	*p = *NewItem()

    if known.ObjectType_ != nil {
        p.ObjectType_ = known.ObjectType_
    }
    if known.Reserved_ != nil {
        p.Reserved_ = known.Reserved_
    }
    if known.UnknownFields_ != nil {
        p.UnknownFields_ = known.UnknownFields_
    }
    if known.Description != nil {
        p.Description = known.Description
    }
    if known.ItemId != nil {
        p.ItemId = known.ItemId
    }
    if known.ItemImageFile != nil {
        p.ItemImageFile = known.ItemImageFile
    }
    if known.ItemName != nil {
        p.ItemName = known.ItemName
    }
    if known.ItemType != nil {
        p.ItemType = known.ItemType
    }
    if known.Location != nil {
        p.Location = known.Location
    }

    // Step 4: Remove known JSON fields from allFields map
	delete(allFields, "$objectType")
	delete(allFields, "$reserved")
	delete(allFields, "$unknownFields")
	delete(allFields, "description")
	delete(allFields, "itemId")
	delete(allFields, "itemImageFile")
	delete(allFields, "itemName")
	delete(allFields, "itemType")
	delete(allFields, "location")

    // Step 5: Assign remaining fields to UnknownFields_
	for key, value := range allFields {
      p.UnknownFields_[key] = value
    }

	return nil
}

func NewItem() *Item {
  p := new(Item)
  p.ObjectType_ = new(string)
  *p.ObjectType_ = "nexus.v4.config.Item"
  p.Reserved_ = map[string]interface{}{"$fv": "v4.r1"}
  p.UnknownFields_ = map[string]interface{}{}



  return p
}



/*
REST response for all response codes in API path /nexus/v4.1/config/items Get operation
*/
type ListItemsApiResponse struct {
  
  ObjectType_ *string `json:"$objectType,omitempty"`
  
  Reserved_ map[string]interface{} `json:"$reserved,omitempty"`
  
  UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
  /*
  
  */
  DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`
  
  Data *OneOfListItemsApiResponseData `json:"data,omitempty"`
  
  Metadata *import2.ApiResponseMetadata `json:"metadata,omitempty"`
}

func (p *ListItemsApiResponse) MarshalJSON() ([]byte, error) {
  // Create Alias to avoid infinite recursion
  type Alias ListItemsApiResponse

  // Step 1: Marshal the known fields
  known, err := json.Marshal(Alias(*p))
  if err != nil {
  	return nil, err
  }

    // Step 2: Convert known to map for merging
    var knownMap map[string]interface{}
    if err := json.Unmarshal(known, &knownMap); err != nil {
    	return nil, err
    }
    delete(knownMap, "$unknownFields")
  
    // Step 3: Merge unknown fields
    for k, v := range p.UnknownFields_ {
    	knownMap[k] = v
    }
  
    // Step 4: Marshal final merged map
    return json.Marshal(knownMap)
}

func (p *ListItemsApiResponse) UnmarshalJSON(b []byte) error {
    // Step 1: Unmarshal into a generic map to capture all fields
    var allFields map[string]interface{}
	if err := json.Unmarshal(b, &allFields); err != nil {
		return err
	}

    // Step 2: Unmarshal into a temporary struct with known fields
	type Alias ListItemsApiResponse
	known := &Alias{}
	if err := json.Unmarshal(b, known); err != nil {
		return err
	}

    // Step 3: Assign known fields
	*p = *NewListItemsApiResponse()

    if known.ObjectType_ != nil {
        p.ObjectType_ = known.ObjectType_
    }
    if known.Reserved_ != nil {
        p.Reserved_ = known.Reserved_
    }
    if known.UnknownFields_ != nil {
        p.UnknownFields_ = known.UnknownFields_
    }
    if known.DataItemDiscriminator_ != nil {
        p.DataItemDiscriminator_ = known.DataItemDiscriminator_
    }
    if known.Data != nil {
        p.Data = known.Data
    }
    if known.Metadata != nil {
        p.Metadata = known.Metadata
    }

    // Step 4: Remove known JSON fields from allFields map
	delete(allFields, "$objectType")
	delete(allFields, "$reserved")
	delete(allFields, "$unknownFields")
	delete(allFields, "$dataItemDiscriminator")
	delete(allFields, "data")
	delete(allFields, "metadata")

    // Step 5: Assign remaining fields to UnknownFields_
	for key, value := range allFields {
      p.UnknownFields_[key] = value
    }

	return nil
}

func NewListItemsApiResponse() *ListItemsApiResponse {
  p := new(ListItemsApiResponse)
  p.ObjectType_ = new(string)
  *p.ObjectType_ = "nexus.v4.config.ListItemsApiResponse"
  p.Reserved_ = map[string]interface{}{"$fv": "v4.r1"}
  p.UnknownFields_ = map[string]interface{}{}



  return p
}

func (p *ListItemsApiResponse) GetData() interface{} {
  if nil == p.Data {
    return nil
  }
  return p.Data.GetValue()
}

func (p *ListItemsApiResponse) SetData(v interface{}) error {
  if nil == p.Data {
    p.Data = NewOneOfListItemsApiResponseData()
  }
  e := p.Data.SetValue(v)
  if nil == e {
    if nil == p.DataItemDiscriminator_ {
      p.DataItemDiscriminator_ = new(string)
    }
    *p.DataItemDiscriminator_ = *p.Data.Discriminator
  }
  return e
}


/*
Geographical location information
*/
type Location struct {
  
  ObjectType_ *string `json:"$objectType,omitempty"`
  
  Reserved_ map[string]interface{} `json:"$reserved,omitempty"`
  
  UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
  /*
  City name
  */
  City *string `json:"city,omitempty"`
  
  Country *Country `json:"country,omitempty"`
  /*
  ZIP or postal code
  */
  Zip *string `json:"zip,omitempty"`
}

func (p *Location) MarshalJSON() ([]byte, error) {
  // Create Alias to avoid infinite recursion
  type Alias Location

  // Step 1: Marshal the known fields
  known, err := json.Marshal(Alias(*p))
  if err != nil {
  	return nil, err
  }

    // Step 2: Convert known to map for merging
    var knownMap map[string]interface{}
    if err := json.Unmarshal(known, &knownMap); err != nil {
    	return nil, err
    }
    delete(knownMap, "$unknownFields")
  
    // Step 3: Merge unknown fields
    for k, v := range p.UnknownFields_ {
    	knownMap[k] = v
    }
  
    // Step 4: Marshal final merged map
    return json.Marshal(knownMap)
}

func (p *Location) UnmarshalJSON(b []byte) error {
    // Step 1: Unmarshal into a generic map to capture all fields
    var allFields map[string]interface{}
	if err := json.Unmarshal(b, &allFields); err != nil {
		return err
	}

    // Step 2: Unmarshal into a temporary struct with known fields
	type Alias Location
	known := &Alias{}
	if err := json.Unmarshal(b, known); err != nil {
		return err
	}

    // Step 3: Assign known fields
	*p = *NewLocation()

    if known.ObjectType_ != nil {
        p.ObjectType_ = known.ObjectType_
    }
    if known.Reserved_ != nil {
        p.Reserved_ = known.Reserved_
    }
    if known.UnknownFields_ != nil {
        p.UnknownFields_ = known.UnknownFields_
    }
    if known.City != nil {
        p.City = known.City
    }
    if known.Country != nil {
        p.Country = known.Country
    }
    if known.Zip != nil {
        p.Zip = known.Zip
    }

    // Step 4: Remove known JSON fields from allFields map
	delete(allFields, "$objectType")
	delete(allFields, "$reserved")
	delete(allFields, "$unknownFields")
	delete(allFields, "city")
	delete(allFields, "country")
	delete(allFields, "zip")

    // Step 5: Assign remaining fields to UnknownFields_
	for key, value := range allFields {
      p.UnknownFields_[key] = value
    }

	return nil
}

func NewLocation() *Location {
  p := new(Location)
  p.ObjectType_ = new(string)
  *p.ObjectType_ = "nexus.v4.config.Location"
  p.Reserved_ = map[string]interface{}{"$fv": "v4.r1"}
  p.UnknownFields_ = map[string]interface{}{}



  return p
}



type OneOfListItemsApiResponseData struct {
  Discriminator *string `json:"-"`
  ObjectType_ *string `json:"-"`
  oneOfType2001 []Item `json:"-"`
  oneOfType400 *import1.ErrorResponse `json:"-"`
}

func NewOneOfListItemsApiResponseData() *OneOfListItemsApiResponseData {
  p := new(OneOfListItemsApiResponseData)
  p.Discriminator = new(string)
  p.ObjectType_ = new(string)
  return p
}

func (p *OneOfListItemsApiResponseData) SetValue (v interface {}) error {
  if nil == p {
    return errors.New(fmt.Sprintf("OneOfListItemsApiResponseData is nil"))
  }
  switch v.(type) {
    case []Item:
      p.oneOfType2001 = v.([]Item)
      if nil == p.Discriminator {p.Discriminator = new(string)}
      *p.Discriminator = "List<nexus.v4.config.Item>"
      if nil == p.ObjectType_ {p.ObjectType_ = new(string)}
      *p.ObjectType_ = "List<nexus.v4.config.Item>"
    case import1.ErrorResponse:
      if nil == p.oneOfType400 {p.oneOfType400 = new(import1.ErrorResponse)}
      *p.oneOfType400 = v.(import1.ErrorResponse)
      if nil == p.Discriminator {p.Discriminator = new(string)}
      *p.Discriminator = *p.oneOfType400.ObjectType_
      if nil == p.ObjectType_ {p.ObjectType_ = new(string)}
      *p.ObjectType_ = *p.oneOfType400.ObjectType_
    default:
      return errors.New(fmt.Sprintf("%T(%v) is not expected type", v,v))
  }
  return nil
}

func (p *OneOfListItemsApiResponseData) GetValue() interface{} {
  if "List<nexus.v4.config.Item>" == *p.Discriminator {
    return p.oneOfType2001
  }
  if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
    return *p.oneOfType400
  }
  return nil
}

func (p *OneOfListItemsApiResponseData) UnmarshalJSON(b []byte) error {
  vOneOfType2001 := new([]Item)
  if err := json.Unmarshal(b, vOneOfType2001); err == nil {
    if len(*vOneOfType2001) == 0 || "nexus.v4.config.Item" == *((*vOneOfType2001)[0].ObjectType_) {
      p.oneOfType2001 = *vOneOfType2001
      if nil == p.Discriminator {p.Discriminator = new(string)}
      *p.Discriminator = "List<nexus.v4.config.Item>"
      if nil == p.ObjectType_ {p.ObjectType_ = new(string)}
      *p.ObjectType_ = "List<nexus.v4.config.Item>"
      return nil
    }
  }
  vOneOfType400 := new(import1.ErrorResponse)
  if err := json.Unmarshal(b, vOneOfType400); err == nil {
    if "nexus.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
      if nil == p.oneOfType400 {p.oneOfType400 = new(import1.ErrorResponse)}
      *p.oneOfType400 = *vOneOfType400
      if nil == p.Discriminator {p.Discriminator = new(string)}
      *p.Discriminator = *p.oneOfType400.ObjectType_
      if nil == p.ObjectType_ {p.ObjectType_ = new(string)}
      *p.ObjectType_ = *p.oneOfType400.ObjectType_
      return nil
    }
  }
  return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfListItemsApiResponseData"))
}

func (p *OneOfListItemsApiResponseData) MarshalJSON() ([]byte, error) {
  if "List<nexus.v4.config.Item>" == *p.Discriminator {
    return json.Marshal(p.oneOfType2001)
  }
  if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
    return json.Marshal(p.oneOfType400)
  }
  return nil, errors.New("No value to marshal for OneOfListItemsApiResponseData")
}


type FileDetail struct {
	Path *string `json:"-"`
	ObjectType_ *string `json:"-"`
}

func NewFileDetail() *FileDetail {
	p := new(FileDetail)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "FileDetail"

	return p
}
