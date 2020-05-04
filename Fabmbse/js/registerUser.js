/*
 * SPDX-License-Identifier: Apache-2.0
 */

'use strict';


//Outside main function because async main cannot run sync methods
const { FileSystemWallet, Gateway, X509WalletMixin } = require('fabric-network');
const path = require('path');
const config = require('./config.json');
const ccpPath = path.resolve(__dirname, '..', '..', 'first-network', 'connection-org1.json');

async function main() {
    try {

        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), config.Common.Wallet);
        const wallet = new FileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        const userExists = await wallet.exists(config.regisUser.userIdentity);
        if (userExists) {
            console.log('An identity for the user "user1" already exists in the wallet');
            return;
        }

        // Check to see if we've already enrolled the admin user.
        const adminExists = await wallet.exists(config.adminEnroll.adminIdentity);
        if (!adminExists) {
            console.log('An identity for the admin user "admin" does not exist in the wallet');
            console.log('Run the enrollAdmin.js application before retrying');
            return;
        }

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccpPath, { wallet, identity: config.adminEnroll.adminIdentity, discovery: { enabled: true, asLocalhost: true } });

        // Get the CA client object from the gateway for interacting with the CA.
        const ca = gateway.getClient().getCertificateAuthority();
        const adminIdentity = gateway.getCurrentIdentity();

        // Register the user, enroll the user, and import the new identity into the wallet.
        const secret = await ca.register({ affiliation: config.regisUser.Affiliation, enrollmentID: config.regisUser.userIdentity, role: config.regisUser.Role }, adminIdentity);
        const enrollment = await ca.enroll({ enrollmentID: config.regisUser.userIdentity, enrollmentSecret: secret });
        const userIdentity = X509WalletMixin.createIdentity(config.Common.OrganisationMSP, enrollment.certificate, enrollment.key.toBytes());
        await wallet.import('user1', userIdentity);
        console.log('Successfully registered and enrolled admin user "user1" and imported it into the wallet');

    } catch (error) {
        console.error(`Failed to register user "user1": ${error}`);
        process.exit(1);
    }
}

main();
