import  hre  from "hardhat";

const { ethers } = await hre.network.getOrCreate();

async function main() {

    const [deployer,admin] = await ethers.getSigners();
    const balance = await ethers.provider.getBalance(deployer.address);
    console.log("Deploying contracts with account:", deployer.address);
    console.log("Account balance:", ethers.formatEther(balance), "ETH");

    if (balance === 0n) {
        throw new Error("Account has no balance. Please fund the account first.");
    }

    const MyVersion_v1 = await ethers.getContractFactory("MyVersion_v1");
    const myVersion_v1 = await MyVersion_v1.deploy();

    console.log("\nDeploying MultiSigWallet contract...");

    await myVersion_v1.waitForDeployment();

    const code = await ethers.provider.getCode(await myVersion_v1.getAddress());
    if (code === "0x") {
        throw new Error("myVersion_v1 deployment failed. Please check the transaction receipt.");
    }

    //部署的逻辑合约
    const impAddressNew = await myVersion_v1.getAddress();

    console.log("myVersion_v1 contract deployed to:",impAddressNew);

     const initData = MyVersion_v1.interface.encodeFunctionData("reinitialize", [2,[deployer.address,admin.address],2]);


    //部署代理合约的admin合约
    const ADMIN_SLOT = "0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103";
    
    console.log("Deploying ProxyAdmin contract...");

      // 4. 部署透明代理，传入逻辑合约、ProxyAdmin、初始化数据
    const TransparentUpgradeableProxy = await ethers.getContractFactory("MultiSigWalletProxy");
    const prox = await TransparentUpgradeableProxy.attach("0x610178dA211FEF7D417bC0e6FeD39F05609AD788");
    
   

    console.log("代理地址:", await prox.getAddress());

     const adminAddress = await ethers.provider.getStorage(await prox.getAddress(), ADMIN_SLOT);
     const adminStorageHex = adminAddress.startsWith("0x") ? adminAddress.slice(2) : adminAddress;
     const actualProxyAdminAddress = ethers.getAddress("0x" + adminStorageHex.slice(-40).padStart(40, '0'));

     console.log("ProxyAdmin contract address:",actualProxyAdminAddress);
      const proxyAdmin = await ethers.getContractAt("ProxyAdmin", actualProxyAdminAddress);
      const adminUser = await proxyAdmin.owner();
      console.log("ProxyAdmin contract owner:",adminUser);
     if(admin.address === adminUser){
         console.log("ProxyAdmin contract already deployed to:",actualProxyAdminAddress);
     }else{
         throw new Error("myVersion_v1 admin合约账户获取失败");
     }
    console.log("ProxyAdmin contract deployed to:",actualProxyAdminAddress);

    await proxyAdmin.connect(admin).upgradeAndCall(await prox.getAddress(), impAddressNew, initData)
   


    console.log("Deployment complete.");
}

// 执行主函数
main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });