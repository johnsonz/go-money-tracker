<html>
<head>
    <title></title>
</head>
<body>
    <form action="/login" method="POST">
    <div>
        <input type="text" name="username" value="{{.Username}}"/>
        <input type="password" name="password"/>
        <input type="submit" value="Login"/>
    </div>
    <div>
        <label style="display:{{.ShowError}}">Username or password is incorrect.</label>
    </div>
</form>
</body>
</html>
