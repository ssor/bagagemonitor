<!DOCTYPE html>
<html>

<head>
    <title>API列表</title>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
    <style type="text/css">
	    body{
		    padding-left: 10px;
		    padding-top: 20px;
	    }
	    .apiName {
	        font-size: 17px;
	        padding-bottom: 10px
	    }
	    .apiContent{
	    	padding-left: 20px;
	    	border-bottom: 5px solid rgba(200,200,200,0.4);
	    	padding-bottom: 10px;
	    	border-top: 1px solid rgba(220,220,222,0.3);
	    	padding-top: 10px;
	    }
	    .apiUrl {
	        background-color: rgba(220, 220, 220, 0.1);
	        margin-top: 10px;
	        padding-top: 10px;
	        padding-bottom: 10px;
	        padding-left: 20px;
	        font-size: 14px;
		    margin-bottom: 10px;
	    }
	    .apiNote{
	    	font-size: 14px;
	    }
	    .apiParaDetail{
		    font-size: 14px;
		    color: rgba(100,100,100,0.8);
	        padding-left: 20px;
		    margin-bottom: 10px;
	    }
    </style>
</head>

<body>
	<div style="    font-size: 30px; padding-bottom: 20px;">API列表</div>
    <div class="apiName">
        添加订单位置信息
    </div>
    <div class="apiParaDetail">添加最新的订单状态信息到系统中，如果新添加的状态信息与现有的订单的状态不同，系统将做出相应处理</div>
    <div class="apiContent">
        <div class="apiNote"> URL: </div>
        <div class="apiUrl">
            POST /addBagage
        </div>

        <div class="apiNote"> 参数: </div>
        <div class="apiUrl">
            {ID:"", Longitude: 39.45, Latitude: 112.45, BagageList: ["001", "002"]}
        </div>
        <div class="apiParaDetail">ID:与订单绑定的打包或者车辆信息;</div>

        <div class="apiNote"> 返回: </div>
        <div class="apiUrl">
            {Code: 0, Message: ""}
        </div>
        <div class="apiParaDetail">Code为0,表示成功;Code为1,表示失败，Message中返回的是出错提示信息 </div>
    </div>
</body>

</html>
