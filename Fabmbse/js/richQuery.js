'use strict';
const { FileSystemWallet, Gateway } = require('fabric-network');
const path = require('path');
const fs = require('fs')
const ccpPath = path.resolve(__dirname, '..', '..', 'first-network', 'connection-org1.json');
const jsonString2 = fs.readFileSync('./selectors.json')

async function main() {
    try {

        
        const walletPath = path.join(process.cwd(), 'wallet');
        const wallet = new FileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

      
        const userExists = await wallet.exists('user1');
        if (!userExists) {
            console.log('An identity for the user "user1" does not exist in the wallet');
            console.log('Run the registerUser.js application before retrying');
            return;
        }

      
        const gateway = new Gateway();
        await gateway.connect(ccpPath, { wallet, identity: 'user1', discovery: { enabled: true, asLocalhost: true } });

       
        const network = await gateway.getNetwork('mychannel');

        var myArgs = process.argv.slice(2);
        //JSON transformation
        var temp = JSON.parse(jsonString2)
        var myJSON = '';
	
	if (myArgs[0] == '0'){
		var temp1 = JSON.stringify(temp.zero)//insert key & value
		myJSON = temp1.replace("{1}", myArgs[1]);
		myJSON = myJSON.replace("{2}", myArgs[2]);
	} else if (myArgs[0] == '1'){
		var temp1 = JSON.stringify(temp.one)//insert value
		myJSON = temp.replace("{1}", myArgs[1]);
	} else if (myArgs[0] == '2'){
		myJSON = JSON.stringify(temp.two);//top level
	} else if (myArgs[0] == '3'){
		myJSON = JSON.stringify(temp.three);//single tier nested
	} else if (myArgs[0] == '4'){
		myJSON = JSON.stringify(temp.four);//not implemented yet
	} else if (myArgs[0] == '5'){
		myJSON = JSON.stringify(temp.five);//query ALL
	} else {
		console.log('Arguments Invalid')
		return
	}
	

        const contract = network.getContract('fabcar');
	   
        const result = await contract.evaluateTransaction('richQueryModel', myJSON);   
        console.log(`Transaction has been evaluated, result is: ${result.toString()}`);

    } catch (error) {
        console.error(`Failed to evaluate transaction: ${error}`);
        process.exit(1);
    }
}



main();

/*

xXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxX
JSON of sample selectors for calling richQuery

selectors below follow the following format:
zero: specify key-value pair
one: specify value
two: top level query
three: nested query
four: double nested query - Not working yet
five: query all

xXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxX

{
"zero":{"selector":{"{1}":"{2}"}}, 
"one":{"selector":{"OrgId":"{1}"}},
"two":{"selector":{"Author": "Joses"}},
"three":{"selector":{"Approval": {"ApprovalStatus":"disapproved"}}},
"four":{"selector":{"Top":{"nest1": {"nest2":"value"}}}},
"five":{"selector":{"everything":{"$exists": false}}},
}

*/
