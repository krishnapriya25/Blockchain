package main

import (
   "bytes"
    "encoding/json"
    "fmt"
    "log"

    "github.com/hyperledger/fabric/core/chaincode/shim"
    sc "github.com/hyperledger/fabric/protos/peer"
)

type Model struct {
    ModelId             string              `json:"ModelId"`
    ModelName           string              `json:"ModelName"`
    OrgId               string              `json:"OrgId"`
    ProjId              string              `json:"ProjId"`
    PrescribedComponentList []PrescribedComponent `json:"PrescribedComponentList"`
}

type PrescribedComponent struct {
    PCId       string `json:"PCId"`
    PCName     string `json:"PCName"`
    CompId     string `json:"CompId"`
    ChildrenPCList []PrescribedComponent `json:"ChildrenPCList"`
}

type Component struct { 
   
    
    CompId string `json:"CompId"` 

    ComponentName string `json:"ComponentName"`   

    ModelId string `json:"ModelId"`  
    
    ParentCompId string `json:"ParentCompId"`

    RelatedToComponents []ToComponent `json:"RelatedToComponents"`
   
    Storage AssetStorage `json:"AssetStorage"`  
  
    Author string `json:"Author"`     
   
    ReleaseTime  string `json:"ReleaseTime"`  
   
    Version string `json:"Version"`

    Subversion string `json:"Subversion"`  
   
    ComponentStatus ComponentStatusInfo `json:"ComponentStatus"`  

    Approval ApprovalInfo `json:"Approval"`              
}

type ToComponent struct {
    
    ToCompId string `json:"ToCompId"`
    
    ToCompType string `json:"ToCompType"`
    
    Description string `json:"Description"`
}

type AssetStorage struct {
   
    StorageMethod string `json:"StorageMethod"`      

    IPFSName string `json:"IPFSName"`                  // IPFS name.

    AssetIPFSAddress string `json:"AssetIPFSAddress"`  // IPFS multihash address

    ISEncrypted bool `json:"IsEncrypted"` 
    
    IPFSFileEncMethod string `json:"IPFSFileEncMethod"`     
   
    IPFSFileEncKey string `json:"IPFSFileEncKey"` 

    AssetRaw string `json:"AssetRaw"`          // component in binary format.

    SourceType string `json:"SourceType"`       
  
    SourceName string  `json:"SourceName"`
    
    SourceFileName string `json:"SourceFileName"`  

    StartTime string `json:"StartTime"`  //time.Time
}

type ComponentStatusInfo struct {

    ComponentStatus string `json:"ComponentStatus"` 

    StatusSince string `json:"StatusSince"`  //time.Time
}

type ApprovalInfo struct {

    ApprovalStatus string `json:"ApprovalStatus"` 
   
    ApprovalTime string `json:"ApprovalTime"`   //time.Time
   
    Approver string `json:"Approver"`
    
    StatusUpdateRequestor string `json:"StatusUpdateRequestor"` 

    StatusUpdateTo string `json:"StatusUpdateTo"`
 
    StatusUpdateReqTime string `json:"StatusUpdateReqTime"`  //time.Time
}

type SmartContract struct {
}


func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
    return shim.Success(nil)
}


func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

    // Retrieve the requested Smart Contract function and arguments
    function, args := APIstub.GetFunctionAndParameters()
    // Route to the appropriate handler function to interact with the ledger appropriately
    if function == "queryModel" {
        return s.queryModel(APIstub, args)
    } else if function == "initLedger" {
        return s.initLedger(APIstub)
    } else if function == "createModel" {
        return s.createModel(APIstub, args)
    } else if function == "queryAllModels" {
        return s.queryAllModels(APIstub)
    } else if function == "createComponent" {
        return s.createComponent(APIstub, args)
    } else if function == "richQueryModel" {
        return s.richQueryModel(APIstub, args)
    }else if function == "queryModelbyId" {
        return s.queryModelbyId(APIstub, args)
    }else if function == "queryComponentbyId" {
        return s.queryComponentbyId(APIstub, args)
    }else if function == "richQueryModel" {
        return s.richQueryModel(APIstub, args)
    } else if function == "modifyComponent" {
        return s.modifyComponent(APIstub, args)
    }

    return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryModel(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

    if len(args) != 1 {
        return shim.Error("Incorrect number of arguments. Expecting 1")
    }

    ModelAsBytes, _ := APIstub.GetState(args[0])
    return shim.Success(ModelAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {

    return shim.Success(nil)
}

func (s *SmartContract) createModel(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

    if len(args) != 1 {
        return shim.Error("Incorrect number of arguments. Expecting 1")
    }

    indexName := "model~id" 
    var result Model
    if err := json.Unmarshal([]byte(args[0]), &result); err != nil {
        log.Fatal(err.Error())
    }

    ModelAsBytes, _ := json.Marshal(result)
        
    attributeIdIndexKey, err := APIstub.CreateCompositeKey(indexName, []string{"model", result.ModelId})
    if err != nil {
        return shim.Error(err.Error())
    }
                        
    APIstub.PutState(attributeIdIndexKey, ModelAsBytes)
    return shim.Success(nil)
}


func ( s *SmartContract) createComponent(APIstub shim.ChaincodeStubInterface, args []string) sc.Response{
    if len(args)!= 1{
        return shim.Error("Expecting number of arguments: 1")
    }
    indexName := "component~id"
    var result Component
    if err := json.Unmarshal([]byte(args[0]), &result); err != nil {
        log.Fatal(err.Error())
    } 
    attributeIdIndexKey, err := APIstub.CreateCompositeKey(indexName, []string{"component", result.CompId})
    if err != nil {
        return shim.Error(err.Error())
    }

    ComponentAsBytes, _ := json.Marshal(result)
    APIstub.PutState(attributeIdIndexKey, ComponentAsBytes)

    return shim.Success(nil)
}

func (s *SmartContract) queryAllModels(APIstub shim.ChaincodeStubInterface) sc.Response {

    startKey := ""
    endKey := ""

    resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
    if err != nil {
        return shim.Error(err.Error())
    }
    defer resultsIterator.Close()

    // buffer is a JSON array containing QueryResults
    var buffer bytes.Buffer
    buffer.WriteString("[")

    bArrayMemberAlreadyWritten := false
    for resultsIterator.HasNext() {
        queryResponse, err := resultsIterator.Next()
        if err != nil {
            return shim.Error(err.Error())
        }
        // Add a comma before array members, suppress it for the first array member
        if bArrayMemberAlreadyWritten == true {
            buffer.WriteString(",")
        }
        buffer.WriteString("{\"Key\":")
        buffer.WriteString("\"")
        buffer.WriteString(queryResponse.Key)
        buffer.WriteString("\"")

        buffer.WriteString(", \"Record\":")
        // Record is a JSON object, so we write as-is
        buffer.WriteString(string(queryResponse.Value))
        buffer.WriteString("}")
        bArrayMemberAlreadyWritten = true
    }
    buffer.WriteString("]")

    fmt.Printf("- queryAllModels:\n%s\n", buffer.String())

    return shim.Success(buffer.Bytes())
}

func (s *SmartContract) queryModelbyId(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

    if len(args) != 1 {
        return shim.Error("Incorrect number of arguments. Expecting 1")
    }
    
    indexName := "model~id"
        
    attrIdIndexKey, err := APIstub.CreateCompositeKey(indexName, []string{"model", args[0]})
    if err != nil {
        return shim.Error(err.Error())
    }

    ModelAsBytes, _ := APIstub.GetState(attrIdIndexKey)

    return shim.Success(ModelAsBytes)
}

func (s *SmartContract) queryComponentbyId(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

    if len(args) != 1 {
        return shim.Error("Incorrect number of arguments. Expecting 1")
    }
    
    indexName := "component~id"
        
    attrIdIndexKey, err := APIstub.CreateCompositeKey(indexName, []string{"component", args[0]})
    if err != nil {
        return shim.Error(err.Error())
    }

    ComponentAsBytes, _ := APIstub.GetState(attrIdIndexKey)

    return shim.Success(ComponentAsBytes)
}

func (s *SmartContract) richQueryModel(stub shim.ChaincodeStubInterface, args []string) sc.Response {
    
    if len(args) != 1  {
        return shim.Error("Incorrect number of arguments. Expecting key and value to query")
    }

    queryString := string(args[0]) 
    fmt.Printf(queryString)
    queryResults, err := getQueryResultForQueryString(stub, queryString)
    if err != nil {
        jsonResp := "{\"Error\":\"Failed to get state for }"
        return shim.Error(jsonResp)
    } else if string(queryResults) == "[]" {
        jsonResp := []byte("{\"Error\":\"Model does not exist: ")       
        return shim.Success(jsonResp)
    }

    return shim.Success(queryResults)
}

func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

    fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

    resultsIterator, err := stub.GetQueryResult(queryString)
    if err != nil {
        return nil, err
    }
    defer resultsIterator.Close()

    buffer, err := constructQueryResponseFromIterator(resultsIterator)
    if err != nil {
        return nil, err
    }
    fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

    return buffer.Bytes(), nil
}

func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {
    
    var buffer bytes.Buffer
    buffer.WriteString("[")
    bArrayMemberAlreadyWritten := false

    for resultsIterator.HasNext() {
        queryResponse, err := resultsIterator.Next()
        if err != nil {
            return nil, err
        }
        
        if bArrayMemberAlreadyWritten == true {
            buffer.WriteString(",")
        }
        buffer.WriteString("{\"Key\":")
        buffer.WriteString("\"")
        buffer.WriteString(queryResponse.Key)
        buffer.WriteString("\"")
        buffer.WriteString(", \"Record\":")
        buffer.WriteString(string(queryResponse.Value))
        buffer.WriteString("}")
        bArrayMemberAlreadyWritten = true
    }
    buffer.WriteString("]")

    return &buffer, nil
}

func (s *SmartContract) modifyComponent(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

    if len(args) != 3 {
        return shim.Error("Incorrect number of arguments. Expecting 3")
    }

    indexName := "component~id"

    attrIdIndexKey, err := APIstub.CreateCompositeKey(indexName, []string{"component", args[0]})
    if err != nil {
        return shim.Error(err.Error())
    }

    ComponentAsBytes, _ := APIstub.GetState(attrIdIndexKey)

    var temp Component
    json.Unmarshal(ComponentAsBytes, &temp)

    
    var switchItem = string(args[1])

     
    switch switchItem{

        case "CompId","ModelId","ParentCompId": 
         return shim.Error(switchItem + "CANNOT BE MODIFIED")
        break

        case "ComponentName":
         temp.ComponentName = args[2]     
        break
        
        case "Author":
         temp.Author = args[2]
        break

        case "ComponentStatus":
         var new_status = args[2]
         if new_status != "in_model"||new_status !="preliminary"||new_status !="depracated"||new_status !="deleted"||new_status !="final"{
            shim.Error("ERROR: Incorrect Status Update !")
            } else{
            temp.ComponentStatus.ComponentStatus = args[2]
            }
        break

        case "ApprovalStatus":
         temp.Approval.ApprovalStatus = args[2] 
        break

        case "Approver":
         temp.Approval.Approver = args[2]
        break

        case "StatusUpdateRequestor":
         temp.Approval.StatusUpdateRequestor = args[2]    
        break

        case "StatusUpdateTo":
         temp.Approval.StatusUpdateTo = args[2]
        break   

        default:
         return shim.Error("Requested Modification cannot be made")
    }

    ComponentAsBytes, _ = json.Marshal(temp)
    APIstub.PutState(attrIdIndexKey, ComponentAsBytes)

    return shim.Success(ComponentAsBytes)
}

func main() {

    // Create a new Smart Contract
    err := shim.Start(new(SmartContract))
    if err != nil {
        fmt.Printf("Error creating new Smart Contract: %s", err)
    }
}