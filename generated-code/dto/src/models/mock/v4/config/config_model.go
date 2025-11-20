/*
 * Generated file models/mock/v4/config/config_model.go.
 *
 * Product version: 1.0.0-SNAPSHOT
 *
 * Part of the GoLang Mock API - REST API for Mock Cat Service
 *
 * (c) 2025 Nutanix Inc.  All rights reserved
 *
 */

/*
  Module mock.v4.config of GoLang Mock API - REST API for Mock Cat Service
*/
package config
import (
  "encoding/json"
)
/*
Cat entity for mock REST API
*/
type Cat struct {
  
  ObjectType_ *string `json:"$objectType,omitempty"`
  
  Reserved_ map[string]interface{} `json:"$reserved,omitempty"`
  
  UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
  /*
  Unique identifier for the cat
  */
  CatId *int `json:"catId,omitempty"`
  /*
  Path to cat image file
  */
  CatImageFile *string `json:"catImageFile,omitempty"`
  /*
  Name of the cat
  */
  CatName *string `json:"catName"`
  /*
  Type of cat
  */
  CatType *string `json:"catType"`
  /*
  Description of the cat
  */
  Description *string `json:"description,omitempty"`
  
  Location *Location `json:"location,omitempty"`
}

func (p *Cat) MarshalJSON() ([]byte, error) {
  type CatProxy Cat

  // Step 1: Marshal known fields via proxy to enforce required fields
  baseStruct := struct {
    *CatProxy
    CatName *string `json:"catName,omitempty"`
    CatType *string `json:"catType,omitempty"`
  }{
    CatProxy : (*CatProxy)(p),
    CatName : p.CatName,
    CatType : p.CatType,
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

func (p *Cat) UnmarshalJSON(b []byte) error {
    // Step 1: Unmarshal into a generic map to capture all fields
    var allFields map[string]interface{}
	if err := json.Unmarshal(b, &allFields); err != nil {
		return err
	}

    // Step 2: Unmarshal into a temporary struct with known fields
	type Alias Cat
	known := &Alias{}
	if err := json.Unmarshal(b, known); err != nil {
		return err
	}

    // Step 3: Assign known fields
	*p = *NewCat()

    if known.ObjectType_ != nil {
        p.ObjectType_ = known.ObjectType_
    }
    if known.Reserved_ != nil {
        p.Reserved_ = known.Reserved_
    }
    if known.UnknownFields_ != nil {
        p.UnknownFields_ = known.UnknownFields_
    }
    if known.CatId != nil {
        p.CatId = known.CatId
    }
    if known.CatImageFile != nil {
        p.CatImageFile = known.CatImageFile
    }
    if known.CatName != nil {
        p.CatName = known.CatName
    }
    if known.CatType != nil {
        p.CatType = known.CatType
    }
    if known.Description != nil {
        p.Description = known.Description
    }
    if known.Location != nil {
        p.Location = known.Location
    }

    // Step 4: Remove known JSON fields from allFields map
	delete(allFields, "$objectType")
	delete(allFields, "$reserved")
	delete(allFields, "$unknownFields")
	delete(allFields, "catId")
	delete(allFields, "catImageFile")
	delete(allFields, "catName")
	delete(allFields, "catType")
	delete(allFields, "description")
	delete(allFields, "location")

    // Step 5: Assign remaining fields to UnknownFields_
	for key, value := range allFields {
      p.UnknownFields_[key] = value
    }

	return nil
}

func NewCat() *Cat {
  p := new(Cat)
  p.ObjectType_ = new(string)
  *p.ObjectType_ = "mock.v4.config.Cat"
  p.Reserved_ = map[string]interface{}{"$fv": "v4.r1"}
  p.UnknownFields_ = map[string]interface{}{}



  return p
}



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
  *p.ObjectType_ = "mock.v4.config.Country"
  p.Reserved_ = map[string]interface{}{"$fv": "v4.r1"}
  p.UnknownFields_ = map[string]interface{}{}



  return p
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
  *p.ObjectType_ = "mock.v4.config.Location"
  p.Reserved_ = map[string]interface{}{"$fv": "v4.r1"}
  p.UnknownFields_ = map[string]interface{}{}



  return p
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
