/** @type import('hardhat/config').HardhatUserConfig */
const accounts = [
  {
    privateKey: "0x79ea8f62d97bc0591a4224c1725fca6b00de5b2cea286fe2e0bb35c5e76be46e",
    balance: "10000000000000000000"
  },
  {
    privateKey: "0x1af2e950272dd403de7a5760d41c6e44d92b6d02797e51810795ff03cc2cda4f",
    balance: "10000000000000000000"
  },
  {
    privateKey: "0xf63d7d8e930bccd74e93cf5662fde2c28fd8be95edb70c73f1bdd863d07f412e",
    balance: "10000000000000000000"
  }
];

module.exports = {
  solidity: "0.8.28",
  networks: {
    hardhat: {
      accounts: accounts,
      chainId: 1337,
      blockGasLimit: 12000000,
      mining: {
        auto: true,
        interval: 5000 // mimics -b 5 (5s block time)
      }
    }
  }
};
