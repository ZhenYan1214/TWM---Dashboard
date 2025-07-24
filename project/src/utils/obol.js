import { ethers } from "ethers";
import  { multicall, wrapCall } from "./multicall";
import { MdArrowForward } from "@vicons/ionicons4";
// Operator Clusters Registry Contract
const ObolOperatorClustersRegistry = "0xaE7B191A31f627b4eB1d4DaC64eaB9976995b433";
// SplitMain Contract
const SplitMain = "0x2ed6c4B5dA6378c7897AC67Ba9e43102Feb694EE";
// wstETH Contract
const WSTETH_CONTRACT = '0x7f39C581F595B53c5cb19bD0b3f8dA6c935E2Ca0';

// Obol Operator Clusters Registry Contract ABI
const ObolOperatorClustersRegistryProxyABI = [{"constant":true,"inputs":[],"name":"proxyType","outputs":[{"name":"proxyTypeId","type":"uint256"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[],"name":"isDepositable","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"implementation","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"appId","outputs":[{"name":"","type":"bytes32"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"kernel","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"inputs":[{"name":"_kernel","type":"address"},{"name":"_appId","type":"bytes32"},{"name":"_initializePayload","type":"bytes"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"payable":true,"stateMutability":"payable","type":"fallback"},{"anonymous":false,"inputs":[{"indexed":false,"name":"sender","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"ProxyDeposit","type":"event"}]
const ObolOperatorClustersRegistryLogicABI = [{"anonymous":false,"inputs":[{"indexed":false,"name":"nodeOperatorId","type":"uint256"},{"indexed":false,"name":"name","type":"string"},{"indexed":false,"name":"rewardAddress","type":"address"},{"indexed":false,"name":"stakingLimit","type":"uint64"}],"name":"NodeOperatorAdded","type":"event"},{"constant":true,"inputs":[],"name":"getActiveNodeOperatorsCount","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"_nodeOperatorId","type":"uint256"},{"name":"_fullInfo","type":"bool"}],"name":"getNodeOperator","outputs":[{"name":"active","type":"bool"},{"name":"name","type":"string"},{"name":"rewardAddress","type":"address"},{"name":"totalVettedValidators","type":"uint64"},{"name":"totalExitedValidators","type":"uint64"},{"name":"totalAddedValidators","type":"uint64"},{"name":"totalDepositedValidators","type":"uint64"}],"payable":false,"stateMutability":"view","type":"function"}]
const ObolOperatorABI = [{"inputs":[{"internalType":"address","name":"_feeRecipient","type":"address"},{"internalType":"uint256","name":"_feeShare","type":"uint256"},{"internalType":"contract ERC20","name":"_stETH","type":"address"},{"internalType":"contract ERC20","name":"_wstETH","type":"address"}],"stateMutability":"nonpayable","type":"constructor"},{"inputs":[],"name":"Invalid_Address","type":"error"},{"inputs":[],"name":"Invalid_FeeRecipient","type":"error"},{"inputs":[{"internalType":"uint256","name":"fee","type":"uint256"}],"name":"Invalid_FeeShare","type":"error"},{"inputs":[],"name":"distribute","outputs":[{"internalType":"uint256","name":"amount","type":"uint256"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"feeRecipient","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"feeShare","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"token","type":"address"}],"name":"rescueFunds","outputs":[{"internalType":"uint256","name":"balance","type":"uint256"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"splitWallet","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"pure","type":"function"},{"inputs":[],"name":"stETH","outputs":[{"internalType":"contract ERC20","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"wstETH","outputs":[{"internalType":"contract ERC20","name":"","type":"address"}],"stateMutability":"view","type":"function"}];
const fullABI = [...ObolOperatorClustersRegistryProxyABI, ...ObolOperatorClustersRegistryLogicABI];
const SplitMainABI = [{"inputs":[],"stateMutability":"nonpayable","type":"constructor"},{"inputs":[],"name":"Create2Error","type":"error"},{"inputs":[],"name":"CreateError","type":"error"},{"inputs":[{"internalType":"address","name":"newController","type":"address"}],"name":"InvalidNewController","type":"error"},{"inputs":[{"internalType":"uint256","name":"accountsLength","type":"uint256"},{"internalType":"uint256","name":"allocationsLength","type":"uint256"}],"name":"InvalidSplit__AccountsAndAllocationsMismatch","type":"error"},{"inputs":[{"internalType":"uint256","name":"index","type":"uint256"}],"name":"InvalidSplit__AccountsOutOfOrder","type":"error"},{"inputs":[{"internalType":"uint256","name":"index","type":"uint256"}],"name":"InvalidSplit__AllocationMustBePositive","type":"error"},{"inputs":[{"internalType":"uint32","name":"allocationsSum","type":"uint32"}],"name":"InvalidSplit__InvalidAllocationsSum","type":"error"},{"inputs":[{"internalType":"uint32","name":"distributorFee","type":"uint32"}],"name":"InvalidSplit__InvalidDistributorFee","type":"error"},{"inputs":[{"internalType":"bytes32","name":"hash","type":"bytes32"}],"name":"InvalidSplit__InvalidHash","type":"error"},{"inputs":[{"internalType":"uint256","name":"accountsLength","type":"uint256"}],"name":"InvalidSplit__TooFewAccounts","type":"error"},{"inputs":[{"internalType":"address","name":"sender","type":"address"}],"name":"Unauthorized","type":"error"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"split","type":"address"}],"name":"CancelControlTransfer","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"split","type":"address"},{"indexed":true,"internalType":"address","name":"previousController","type":"address"},{"indexed":true,"internalType":"address","name":"newController","type":"address"}],"name":"ControlTransfer","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"split","type":"address"}],"name":"CreateSplit","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"split","type":"address"},{"indexed":true,"internalType":"contract ERC20","name":"token","type":"address"},{"indexed":false,"internalType":"uint256","name":"amount","type":"uint256"},{"indexed":true,"internalType":"address","name":"distributorAddress","type":"address"}],"name":"DistributeERC20","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"split","type":"address"},{"indexed":false,"internalType":"uint256","name":"amount","type":"uint256"},{"indexed":true,"internalType":"address","name":"distributorAddress","type":"address"}],"name":"DistributeETH","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"split","type":"address"},{"indexed":true,"internalType":"address","name":"newPotentialController","type":"address"}],"name":"InitiateControlTransfer","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"split","type":"address"}],"name":"UpdateSplit","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"account","type":"address"},{"indexed":false,"internalType":"uint256","name":"ethAmount","type":"uint256"},{"indexed":false,"internalType":"contract ERC20[]","name":"tokens","type":"address[]"},{"indexed":false,"internalType":"uint256[]","name":"tokenAmounts","type":"uint256[]"}],"name":"Withdrawal","type":"event"},{"inputs":[],"name":"PERCENTAGE_SCALE","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"split","type":"address"}],"name":"acceptControl","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"split","type":"address"}],"name":"cancelControlTransfer","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address[]","name":"accounts","type":"address[]"},{"internalType":"uint32[]","name":"percentAllocations","type":"uint32[]"},{"internalType":"uint32","name":"distributorFee","type":"uint32"},{"internalType":"address","name":"controller","type":"address"}],"name":"createSplit","outputs":[{"internalType":"address","name":"split","type":"address"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"split","type":"address"},{"internalType":"contract ERC20","name":"token","type":"address"},{"internalType":"address[]","name":"accounts","type":"address[]"},{"internalType":"uint32[]","name":"percentAllocations","type":"uint32[]"},{"internalType":"uint32","name":"distributorFee","type":"uint32"},{"internalType":"address","name":"distributorAddress","type":"address"}],"name":"distributeERC20","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"split","type":"address"},{"internalType":"address[]","name":"accounts","type":"address[]"},{"internalType":"uint32[]","name":"percentAllocations","type":"uint32[]"},{"internalType":"uint32","name":"distributorFee","type":"uint32"},{"internalType":"address","name":"distributorAddress","type":"address"}],"name":"distributeETH","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"split","type":"address"}],"name":"getController","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"account","type":"address"},{"internalType":"contract ERC20","name":"token","type":"address"}],"name":"getERC20Balance","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"account","type":"address"}],"name":"getETHBalance","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"split","type":"address"}],"name":"getHash","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"split","type":"address"}],"name":"getNewPotentialController","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"split","type":"address"}],"name":"makeSplitImmutable","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address[]","name":"accounts","type":"address[]"},{"internalType":"uint32[]","name":"percentAllocations","type":"uint32[]"},{"internalType":"uint32","name":"distributorFee","type":"uint32"}],"name":"predictImmutableSplitAddress","outputs":[{"internalType":"address","name":"split","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"split","type":"address"},{"internalType":"address","name":"newController","type":"address"}],"name":"transferControl","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"split","type":"address"},{"internalType":"contract ERC20","name":"token","type":"address"},{"internalType":"address[]","name":"accounts","type":"address[]"},{"internalType":"uint32[]","name":"percentAllocations","type":"uint32[]"},{"internalType":"uint32","name":"distributorFee","type":"uint32"},{"internalType":"address","name":"distributorAddress","type":"address"}],"name":"updateAndDistributeERC20","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"split","type":"address"},{"internalType":"address[]","name":"accounts","type":"address[]"},{"internalType":"uint32[]","name":"percentAllocations","type":"uint32[]"},{"internalType":"uint32","name":"distributorFee","type":"uint32"},{"internalType":"address","name":"distributorAddress","type":"address"}],"name":"updateAndDistributeETH","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"split","type":"address"},{"internalType":"address[]","name":"accounts","type":"address[]"},{"internalType":"uint32[]","name":"percentAllocations","type":"uint32[]"},{"internalType":"uint32","name":"distributorFee","type":"uint32"}],"name":"updateSplit","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"walletImplementation","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"account","type":"address"},{"internalType":"uint256","name":"withdrawETH","type":"uint256"},{"internalType":"contract ERC20[]","name":"tokens","type":"address[]"}],"name":"withdraw","outputs":[],"stateMutability":"nonpayable","type":"function"},{"stateMutability":"payable","type":"receive"}]
const SplitMainInterface = new ethers.Interface(SplitMainABI);


// 從環境變數獲取配置，如果沒有則使用預設值
const EtherscanAPIKey = import.meta.env.VITE_ETHERSCAN_API_KEY || "RKFRN7F2EGCMFWPMVTX5RE4ZIQ7E1RHJXI";
const RPC_URL = import.meta.env.VITE_RPC_URL || "https://convincing-divine-slug.quiknode.pro/9e835fa839e702f4140f3a80b65672cfb3d5950b/";
const SECONDS_PER_BLOCK = Number(import.meta.env.VITE_SECONDS_PER_BLOCK) || 12;              // 主網平均
const BLOCKS_PER_DAY    = 86_400 / SECONDS_PER_BLOCK;

// 配置載入日誌（僅在開發模式下顯示）
if (import.meta.env.DEV) {
  console.log('🔧 Obol 工具配置:', {
    SECONDS_PER_BLOCK,
    BATCH_SIZE: Number(import.meta.env.VITE_BATCH_SIZE) || 8,
    BATCH_DELAY: Number(import.meta.env.VITE_BATCH_DELAY) || 1000,
    CHART_BATCH_SIZE: Number(import.meta.env.VITE_CHART_BATCH_SIZE) || 6,
    CHART_BATCH_DELAY: Number(import.meta.env.VITE_CHART_BATCH_DELAY) || 1000,
    RPC_URL: RPC_URL.substring(0, 50) + '...',
    ETHERSCAN_API_KEY: EtherscanAPIKey ? '✅ 已設置' : '❌ 未設置'
  });
}

const provider = new ethers.JsonRpcProvider(RPC_URL);
const etherscanProvider = new ethers.EtherscanProvider('homestead',EtherscanAPIKey);
const obolOperatorClustersRegistryContract = new ethers.Contract(ObolOperatorClustersRegistry, fullABI, provider);
const splitMainContract = new ethers.Contract(SplitMain, SplitMainABI, provider);

// 請求節流配置 - 從環境變數獲取
const BATCH_SIZE = Number(import.meta.env.VITE_BATCH_SIZE) || 10; // 每批最多8個請求
const BATCH_DELAY = Number(import.meta.env.VITE_BATCH_DELAY) || 1000; // 批次間延遲1秒

// 圖表數據專用的併發配置（適應 QuickNode 15 req/sec 限制）
const CHART_BATCH_SIZE = Number(import.meta.env.VITE_CHART_BATCH_SIZE) || 6; // 圖表數據每批最多6個區塊
const CHART_BATCH_DELAY = Number(import.meta.env.VITE_CHART_BATCH_DELAY) || 1000; // 圖表數據批次間延遲1秒

// 延遲函數
const delay = (ms) => new Promise(resolve => setTimeout(resolve, ms));



// 圖表數據專用的批次處理（更快的併發）
const processChartBatches = async (items, batchProcessor) => {
  const results = [];
  const totalBatches = Math.ceil(items.length / CHART_BATCH_SIZE);
  
  for (let i = 0; i < items.length; i += CHART_BATCH_SIZE) {
    const batch = items.slice(i, i + CHART_BATCH_SIZE);
    const batchNumber = Math.floor(i / CHART_BATCH_SIZE) + 1;
    const requestCount = batch.length * 2; // 每個區塊需要 2 個請求
    console.log(`📊 圖表數據批次 ${batchNumber}/${totalBatches}: 併發處理 ${batch.length} 個區塊 (${requestCount} 個請求)`);
    
    try {
      const batchResults = await batchProcessor(batch);
      results.push(...batchResults);
      console.log(`✅ 批次 ${batchNumber} 完成，成功處理 ${batchResults.filter(r => r !== null).length}/${batch.length} 個區塊`);
    } catch (error) {
      console.error(`❌ 圖表數據批次 ${batchNumber} 失敗:`, error);
      // 添加空結果以保持索引一致
      results.push(...new Array(batch.length).fill(null));
    }
    
    // 如果不是最後一批，等待較短時間
    if (i + CHART_BATCH_SIZE < items.length) {
      console.log(`⏳ 頻率控制: 等待 ${CHART_BATCH_DELAY}ms 後處理下一批次 (保持 ≤14 req/sec)...`);
      await delay(CHART_BATCH_DELAY);
    }
  }
  
  return results;
};
  
export const ether_obol = {
  getObolOperatorClustersRegistry: async () => {
    try {
      console.log('Fetching active node operators count...');
      const activeNodeOperators = await obolOperatorClustersRegistryContract.getActiveNodeOperatorsCount()
      const operatorCount = Number(activeNodeOperators) // 轉換 BigNumber 為普通數字
      console.log('Active node operators count:', operatorCount)
      
      if (operatorCount === 0) {
        console.log('No active node operators found')
        return { activeNodeOperators: 0, validOperators: [] }
      }
      
      // 創建操作者索引數組
      const operatorIndices = Array.from({ length: operatorCount }, (_, i) => i);
      const RegistryCalls = operatorIndices.map(index => wrapCall(ObolOperatorClustersRegistry, fullABI, "getNodeOperator", [index, true]));
      const results = await multicall(RegistryCalls, provider);
      console.log(results);
      return { activeNodeOperators, validOperators: results };
    } catch (error) {
      console.error('Error in getObolOperatorClustersRegistry:', error);
      throw error;
    }
  },
  getObolOperatorHistoryValidators: async (id, period) => {
    const startTime = Date.now();
    console.log(`🚀 開始載入操作者 ${id} 的 ${period} 歷史數據`);

    const PERIOD_CONFIG = {
      '1m': {days: 30, stepDays: 3},
      '1y': {days: 365, stepDays: 30},
      'all': {days: 1000, stepDays: 60} // 'all' 設為 2 年，每 30 天一個數據點
    };

    const { days, stepDays } = PERIOD_CONFIG[period];
    const latestBlock = await provider.getBlockNumber();
    const latestBlkNum = BigInt(latestBlock);
    
    let startBlkNum = latestBlkNum - BigInt(Math.floor(days * BLOCKS_PER_DAY));
    let stepBlocks = BigInt(Math.floor(stepDays * BLOCKS_PER_DAY));
    
    if (period === 'all') {
      try {
        const eventTopic = obolOperatorClustersRegistryContract.getEvent('NodeOperatorAdded(uint256,string,address,uint64)').getFragment().topicHash;
        const hexId = ethers.toBeHex(id,32)
        console.log('✅ 註冊事件主題:', eventTopic)
        const eventURL = `https://api.etherscan.io/v2/api?chainid=1&module=logs&action=getLogs&fromBlock=${startBlkNum}&toBlock=${latestBlkNum}&topic0=${eventTopic}&page=1&offset=1000&apikey=${EtherscanAPIKey}`
        const response = await fetch(eventURL);
        const data = await response.json();
        console.log('✅ 註冊事件:', data)
        if (data.result.length > 0) {
          console.log('✅ 註冊事件:', data.result)
          for (const event of data.result) {
            if (event.data.includes(hexId)) {
              startBlkNum = BigInt(event.blockNumber);
              const timestamp = event.timeStamp;
              const date = new Date(timestamp * 1000).toISOString().split('T')[0];
              console.log('✅ 註冊事件:', startBlkNum, date)
            }
          }
        } else {
          console.log('❌ 註冊事件不存在')
        }
      } catch (error) {
        console.error('Error in getObolOperatorHistoryValidators:', error);
        throw error;
      }
    }
    // 生成所有需要查詢的區塊號
    const blocks = [];
    for (let b = startBlkNum; b <= latestBlkNum; b += stepBlocks) {
      blocks.push(Number(b));
    }
    
    const estimatedRequests = blocks.length * 2; // 每個區塊需要 2 個請求
    const estimatedTime = Math.ceil(blocks.length / CHART_BATCH_SIZE) * (CHART_BATCH_DELAY / 1000);
    console.log(`📋 需要查詢 ${blocks.length} 個區塊的數據，預估時間段範圍: ${Math.round(days)} 天`);
    console.log(`🚦 請求配置: 每批 ${CHART_BATCH_SIZE} 個區塊 (${CHART_BATCH_SIZE * 2} 個請求), 批次間隔 ${CHART_BATCH_DELAY}ms`);
    console.log(`⏱️ 預估總請求數: ${estimatedRequests} 個，預估完成時間: ${estimatedTime} 秒`);

    // 批次處理器函數 - 併發處理區塊數據
    const batchProcessor = async (blockBatch) => {
      const promises = blockBatch.map(async (blockNumber) => {
        const maxRetries = 2;
        let lastError;
        
        for (let retry = 0; retry <= maxRetries; retry++) {
          try {
            // 併發執行兩個請求：操作者數據和區塊信息
            const [validator, blockInfo] = await Promise.all([
              obolOperatorClustersRegistryContract.getNodeOperator(id, true, {blockTag: blockNumber}),
              provider.getBlock(blockNumber)
            ]);

            return {
              blockNumber: blockNumber,
              timestamp: blockInfo.timestamp * 1000, // 轉換為毫秒
              date: new Date(blockInfo.timestamp * 1000).toISOString().split('T')[0],
              data: {
                active: validator[0],
                name: validator[1],
                rewardAddress: validator[2],
                totalVettedValidators: Number(validator[3]),
                totalExitedValidators: Number(validator[4]),
                totalAddedValidators: Number(validator[5]),
                totalDepositedValidators: Number(validator[6])
              }
            };
                      } catch (error) {
              lastError = error;
              if (retry < maxRetries) {
                // 檢查是否為頻率限制錯誤
                const isRateLimitError = error.message?.includes('request limit reached') || 
                                       error.code === 'UNKNOWN_ERROR' && error.error?.code === -32007;
                
                // 如果是頻率限制錯誤，延遲更長時間
                const retryDelay = isRateLimitError ? 2000 + (retry * 1000) : 500 * (retry + 1);
                
                await delay(retryDelay);
                console.warn(`🔄 區塊 ${blockNumber} 第 ${retry + 1} 次重試 (延遲 ${retryDelay}ms)...`);
              }
            }
        }
        
        console.warn(`❌ 區塊 ${blockNumber} 載入失敗 (${maxRetries + 1} 次嘗試):`, lastError?.message || 'Unknown error');
        return null; // 返回 null 表示失敗
      });

      return Promise.all(promises);
    };

    // 使用圖表專用的批次處理來併發載入數據（更快的併發處理）
    const results = await processChartBatches(blocks, batchProcessor);
    
    // 過濾掉失敗的請求並按時間戳排序
    const validators = results
      .filter(result => result !== null)
      .sort((a, b) => a.timestamp - b.timestamp);
    
    const endTime = Date.now();
    const duration = (endTime - startTime) / 1000;
    const successRate = ((validators.length / blocks.length) * 100).toFixed(1);
    const actualRequestCount = blocks.length * 2; // 總請求數
    const actualRequestRate = (actualRequestCount / duration).toFixed(1);
    
    console.log(`🎉 載入完成！成功 ${validators.length}/${blocks.length} 個數據點 (${successRate}%)`);
    console.log(`⚡ 總耗時: ${duration.toFixed(2)} 秒，平均每個區塊: ${(duration / blocks.length).toFixed(3)} 秒`);
    console.log(`🚦 實際請求頻率: ${actualRequestRate} req/sec (限制: 15 req/sec)`);
    
    return validators;
  },
  getObolOperatorSplitWallets: async (operatorAddress) => {
    const operatorContract = new ethers.Contract(operatorAddress, ObolOperatorABI, provider);
    const splitWallet = await operatorContract.splitWallet();
    return splitWallet;
  },
  getObolOperatorRewardshare: async (splitWallet) => {
    const url = `https://api.etherscan.io/v2/api?chainid=1&module=account&action=txlistinternal&address=${splitWallet}&startblock=0&endblock=99999999&sort=desc&apikey=${EtherscanAPIKey}`;
    const transactions = await fetch(url).then(res => res.json()).then(data => data.result);
    const createTx = transactions.find(tx => tx.type == "create2")
    const createOriginTx = await etherscanProvider.getTransaction(createTx.hash);
    const decodedData = SplitMainInterface.decodeFunctionData('createSplit', createOriginTx.data);
    if (decodedData) {
      const rewardAddress = decodedData[0];
      const rewardShare = decodedData[1];

      return {rewardAddress,rewardShare};
    } else {
      throw new Error('Failed to decode createSplit function data');
    }
  },
  getObolOperatorClaimableReward: async (rewardAddresses) => {
    if(rewardAddresses.length >= 1){
      // 索引陣列
      const rewardAddressIndices = Array.from({length: rewardAddresses.length}, (_, i) => i);
      // 透過索引陣列去映射出對應的Call格式
      const rewardCalls = rewardAddressIndices.map(index => wrapCall(SplitMain, SplitMainABI, "getERC20Balance", [rewardAddresses[index],WSTETH_CONTRACT]));
      const results = await multicall(rewardCalls, provider);
      return results;
    }else{
      throw new Error('No reward addresses provided');
    }
  },



  getObolOperatorTokenTx: async (splitWallet) => {
    const url = `https://api.etherscan.io/v2/api?chainid=1&module=account&action=tokentx&address=${splitWallet}&startblock=0&endblock=99999999&sort=desc&apikey=${EtherscanAPIKey}`;

    const transactions = await fetch(url).then(res => res.json()).then(data => data.result);
    const tokenTx = transactions.find(tx => tx.input.startsWith("0xe4fc6b6  d"));
    return tokenTx;
  },

  // 獲取 wstETH 代幣詳細交易記錄
  getObolOperatorWstETHTransactions: async (splitWallet) => {
    try {
      console.log('🔍 開始載入 wstETH 交易記錄:', splitWallet);
      
      const url = `https://api.etherscan.io/api?module=account&action=tokentx&address=${splitWallet}&startblock=0&endblock=99999999&sort=desc&apikey=${EtherscanAPIKey}`;
      const response = await fetch(url);
      const data = await response.json();
      
      if (data.status !== '1') {
        throw new Error(`Etherscan API 錯誤: ${data.message}`);
      }
      
      const transactions = data.result || [];
      
      
      
      // 過濾 wstETH 交易
      const wstETHTransactions = transactions.filter(tx => 
        tx.contractAddress && tx.contractAddress.toLowerCase() === WSTETH_CONTRACT.toLowerCase()
      );
      
      console.log(`📊 找到 ${wstETHTransactions.length} 筆 wstETH 交易`);
      
      // 處理交易數據，區分收入和支出
      const processedTransactions = wstETHTransactions.map(tx => {
        const value = parseFloat(tx.value) / Math.pow(10, parseInt(tx.tokenDecimal || 18));
        const isIncoming = tx.to.toLowerCase() === splitWallet.toLowerCase();
        
        return {
          hash: tx.hash,
          timeStamp: parseInt(tx.timeStamp),
          from: tx.from,
          to: tx.to,
          value: value,
          valueWei: tx.value,
          isIncoming: isIncoming,
          blockNumber: parseInt(tx.blockNumber),
          date: new Date(parseInt(tx.timeStamp) * 1000).toISOString().split('T')[0],
          tokenName: tx.tokenName,
          tokenSymbol: tx.tokenSymbol,
          tokenDecimal: tx.tokenDecimal
        };
      });
      
      return processedTransactions.sort((a, b) => b.timeStamp - a.timeStamp);
      
    } catch (error) {
      console.error('❌ 載入 wstETH 交易記錄失敗:', error);
      throw error;
    }
  },

  // 統計 wstETH 數據
  getObolOperatorWstETHSummary: async (splitWallet) => {
    try {
      console.log('📈 開始統計 wstETH 數據:', splitWallet);
      
      const transactions = await ether_obol.getObolOperatorWstETHTransactions(splitWallet);
      
      // 計算總收入、總支出和餘額
      let totalReceived = 0;
      let totalDistributed = 0;
      let transactionCount = 0;
      
      const monthlyData = {};
      
      transactions.forEach(tx => {
        if (tx.isIncoming) {
          totalReceived += tx.value;
        } else {
          totalDistributed += tx.value;
        }
        transactionCount++;
        
        // 按月份統計
        const monthKey = tx.date.substring(0, 7); // YYYY-MM 格式
        if (!monthlyData[monthKey]) {
          monthlyData[monthKey] = {
            received: 0,
            distributed: 0,
            transactions: 0
          };
        }
        
        if (tx.isIncoming) {
          monthlyData[monthKey].received += tx.value;
        } else {
          monthlyData[monthKey].distributed += tx.value;
        }
        monthlyData[monthKey].transactions++;
      });
      
      const currentBalance = totalReceived - totalDistributed;
      
      // 轉換月度數據為圖表格式
      const monthlyChartData = Object.keys(monthlyData)
        .sort()
        .map(month => ({
          month: month,
          date: new Date(month + '-01'),
          received: monthlyData[month].received,
          distributed: monthlyData[month].distributed,
          net: monthlyData[month].received - monthlyData[month].distributed,
          transactions: monthlyData[month].transactions
        }));
      
      console.log(`✅ wstETH 數據統計完成: 總收入 ${totalReceived.toFixed(6)} wstETH`);
      
      return {
        totalReceived: totalReceived,
        totalDistributed: totalDistributed,
        currentBalance: currentBalance,
        transactionCount: transactionCount,
        transactions: transactions,
        monthlyData: monthlyChartData,
        lastUpdated: new Date().toISOString(),
        
        // 預留收益預估數據結構（邏輯留空）
        estimatedEarnings: {
          daily: null,      // 預估日收益
          monthly: null,    // 預估月收益
          yearly: null,     // 預估年收益
          apy: null,        // 年化收益率
          nextDistribution: null, // 下次分配預估時間
          // 此處可由用戶自行實現收益預估邏輯
        }
      };
      
    } catch (error) {
      console.error('❌ 統計 wstETH 數據失敗:', error);
      throw error;
    }
  },

  getPridictionWstETH: async (lidoapr, activeValidators, type, isComputable, latestTimestamp) => {
    const share = (type == "Obol" ? 0.07 : 0.08)
    if(isComputable){
      const revenue = lidoapr * activeValidators * share * 32 
      const revenuePerSecond = revenue / 365 / 24 / 60 / 60
      const diffSeconds = (Date.now()/1000) - latestTimestamp
      const result = revenuePerSecond * diffSeconds
      return result
    }
  },

  getLidoProtocolAPR: async () => {
    try {
      const baseUrl = import.meta.env.VITE_LIDO_API_BASE_URL || 'https://eth-api.lido.fi';
      const url = `${baseUrl}/v1/protocol/steth/apr/sma`;
      const response = await fetch(url);
      const data = await response.json();
      console.log('Lido protocol APR 原始數據:', data);
      
      // 解析 APR 數據
      let apr = null;
      if (data.data.aprs && Array.isArray(data.data.aprs) && data.data.aprs.length > 0) {
        const latestAPR = data.data.aprs[data.data.aprs.length - 1];
        if (latestAPR && latestAPR.apr !== undefined) {
          apr = Number(latestAPR.apr);
          console.log('✅ 從 aprs 陣列獲取最新 APR:', apr);
        }
      }
      
      
      // 轉換為小數形式 (API 返回的可能是百分比形式)
      if (apr !== null) {
        // 如果 API 返回的是百分比形式 (例如 2.672)，轉換為小數形式 (0.02672)
        const finalAPR = apr > 1 ? apr / 100 : apr;
        console.log('✅ 最終 APR (小數形式):', finalAPR);
        return finalAPR;
      }
      
      throw new Error('無法從 API 回應中解析 APR 數據');
      
    } catch (error) {
      console.error('❌ 獲取 Lido Protocol APR 失敗:', error);
      throw error;
    }
  }


}