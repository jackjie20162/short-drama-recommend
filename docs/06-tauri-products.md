# Tauri 产品端

admin-tauri：运营/内容后台，登录入口、短剧发布、短剧管理、上下架、免费/付费标记。
client-tauri：用户客户端，推荐 Feed、短剧详情、免费/付费展示、购买入口。

后台启动：
cd admin-tauri
npm install
npm run tauri dev

客户端启动：
cd client-tauri
npm install
npm run tauri dev

后台 API：8081
客户端 API：8080

当前 paid 字段已经进入后台发布和客户端展示；真正的支付订单、支付回调、支付状态仍由独立 order/payment 服务完成，客户端不能直接把购买标记为成功。