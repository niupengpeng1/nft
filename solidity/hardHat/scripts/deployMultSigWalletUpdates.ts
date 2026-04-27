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

    const MultiSigWallet = await ethers.getContractFactory("MultiSigWallet");
    const multiSigWallet = await MultiSigWallet.deploy();

    console.log("\nDeploying MultiSigWallet contract...");

    await multiSigWallet.waitForDeployment();

    const code = await ethers.provider.getCode(await multiSigWallet.getAddress());
    if (code === "0x") {
        throw new Error("multiSigWallet deployment failed. Please check the transaction receipt.");
    }

    //部署的逻辑合约
    const impAddress = await multiSigWallet.getAddress();

    console.log("MultiSigWallet contract deployed to:",impAddress);

     const initData = MultiSigWallet.interface.encodeFunctionData("initialize", [[deployer.address,admin.address],2]);


    //部署代理合约的admin合约
    const ADMIN_SLOT = "0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103";
    
    console.log("Deploying ProxyAdmin contract...");

      // 4. 部署透明代理，传入逻辑合约、ProxyAdmin、初始化数据
    const TransparentUpgradeableProxy = await ethers.getContractFactory("MultiSigWalletProxy");
    const proxy = await TransparentUpgradeableProxy.deploy(
        impAddress,
        admin,
        initData               // 这里就是“构造函数”数据
    );
    await proxy.waitForDeployment();

    console.log("代理地址:", await proxy.getAddress());

     const adminAddress = await ethers.provider.getStorage(await proxy.getAddress(), ADMIN_SLOT);
     const adminStorageHex = adminAddress.startsWith("0x") ? adminAddress.slice(2) : adminAddress;
     const actualProxyAdminAddress = ethers.getAddress("0x" + adminStorageHex.slice(-40).padStart(40, '0'));
    console.log("ProxyAdmin contract deployed to:",actualProxyAdminAddress);

    console.log("Deployment complete.");
}

// 执行主函数
main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });