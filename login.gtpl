<html>
<head>
<title>模拟微信网页登录</title>
</head>
<body>
    <div class="login_box" ng-class="{hide: isScan}">
        <div class="qrcode">
            <img class="img" mm-src="{{qrcodeUrl}}" src="res/blank.gif"/>
            <p class="sub_title">扫描二维码登录</p>
        </div>
    </div>
    <div class="avatar" ng-class="{show: isScan}">
        <img class="img" mm-src="{{userAvatar}}" src="res/blank.gif"/>
        <h4 class="sub_title">扫描成功</h4>
        <p class="tips">请在手机上点击确认以登录</p>
    </div>
</body>
</html>

