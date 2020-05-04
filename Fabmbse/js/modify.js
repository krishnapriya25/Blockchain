'use strict';

const { FileSystemWallet, Gateway } = require('fabric-network');
const path = require('path');
const fs = require('fs')
const ccpPath = path.resolve(__dirname, '..', '..', 'first-network', 'connection-org1.json');
//const jsonString = fs.readFileSync('./test.json')  //specify the file path


async function main() {
    try {

        const walletPath = path.join(process.cwd(), 'wallet');
        const wallet = new FileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

       
        const userExists = await wallet.exists('user1');
        if (!userExists) {
            console.log('An identity for the user "admin" does not exist in the wallet');
            console.log('Run the registerUser.js application before retrying');
            return;
        }

 
        const gateway = new Gateway();
        await gateway.connect(ccpPath, { wallet, identity: 'user1', discovery: { enabled: true, asLocalhost: true } });

       
        const network = await gateway.getNetwork('mychannel');

       
        const contract = network.getContract('fabcar');

        
	    var result = await contract.submitTransaction('modifyComponent', process.argv[2], process.argv[3], process.argv[4]);
	
        console.log('Transaction has been submitted'); 
	    console.log(`Query result is: ${result.toString()}`);
       

        await gateway.disconnect();

    } catch (error) {
        console.error(`Failed to submit transaction: ${error}`);
        process.exit(1);
    }
}

main();
