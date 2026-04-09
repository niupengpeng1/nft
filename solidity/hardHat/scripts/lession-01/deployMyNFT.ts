import { network } from "hardhat";

const { ethers } = await network.connect();


async function main() {

    const [deployer] = await ethers.getSigners();
    const balance = await ethers.provider.getBalance(deployer.address);
    
    console.log("Account balance:", ethers.formatEther(balance), "ETH");

    if (balance === 0n) {
        throw new Error("Account has no balance. Please fund the account first.");
    }

    const myNFT = await ethers.deployContract("MyNFT");

    console.log("\nDeploying myNFT contract...");

    await myNFT.waitForDeployment();

    console.log("Deploying contracts with myNFT:", await myNFT.getAddress());
    const code = await ethers.provider.getCode(await myNFT.getAddress());
    if (code === "0x") {
        throw new Error("myNFT deployment failed. Please check the transaction receipt.");
    }


    console.log("myNFT contract deployed to:", await myNFT.getAddress());

    const tokenID = await myNFT.mint("123",{value:ethers.parseEther("0.01")});
    
    const receipt = await tokenID.wait();

    // 从事件里拿 tokenId ✅ 唯一正确方式
    if(receipt){
        const event = receipt.logs[0];
        // const tokenId = event.args[1]
        console.log("-----",event)
    }
    

    console.log("mint return tokenID  :", tokenID.toString());

    console.log("Deployment complete.");
}


// 执行主函数
main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });
