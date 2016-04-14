<html>

<head>
    <title>Login</title>

    <link rel="stylesheet" href="static/bootstrap-3.3.6/css/bootstrap.min.css">
    <link rel="stylesheet" href="static/bootstrap-3.3.6/css/bootstrap-theme.min.css">
    <link rel="stylesheet" href="static/css/style.css">
</head>

<body>
    <div class="wrap-form">
        <form action="/login" method="POST" class="form-horizontal">
            <div class="form-group">
                <label for="inputusername" class="col-sm-2 control-label">Username</label>
                <div class="col-sm-10">
                    <input type="text" class="form-control" value="{{.Username}}" name="username" id="inputusername" placeholder="Username">
                </div>
            </div>
            <div class="form-group {{.HasError}}">
                <label for="inputPassword3" class="col-sm-2 control-label">Password</label>
                <div class="col-sm-10">
                    <input type="password" class="form-control" name="password" id="inputPassword3" placeholder="Password">
                </div>
            </div>
            <div class="form-group">
                <div class="col-sm-offset-2 col-sm-10">
                    <button type="submit" class="btn btn-primary">Sign in</button>
                </div>
            </div>
        </form>
    </div>
    <script src="/static/js/jquery-1.12.3.min.js"></script>
    <script src="static/bootstrap-3.3.6/js/bootstrap.min.js"></script>
</body>

</html>
