<html>

<head>
    <title>Detail</title>
    <link rel="stylesheet" href="/static/jquery-ui-1.11.4/jquery-ui.min.css">
    <link rel="stylesheet" href="static/bootstrap-3.3.6/css/bootstrap.min.css">
    <link rel="stylesheet" href="static/bootstrap-3.3.6/css/bootstrap-theme.min.css">
    <link rel="stylesheet" href="static/css/style.css">
</head>

<body>
    <form action="/detail" method="POST" enctype="multipart/form-data">
        <input type="hidden" name="itemid" value="{{.ItemID}}" readonly="readonly" />
        <div class="wrap-field">
            <table class="table table-condensed">
                <tr>
                    <td>Name:</td>
                    <td>
                        <input type="text" name="name" />
                    </td>
                </tr>
                <tr>
                    <td>Price:</td>
                    <td>
                        <input type="number" step="0.01" min="0" name="price" />
                    </td>
                </tr>
                <tr>
                    <td>Quantity:</td>
                    <td>
                        <input type="number" min="0" name="quantity" value="1" />
                    </td>
                </tr>
                <tr>
                    <td>Label One:</td>
                    <td>
                        <input type="file" name="labelone" />
                    </td>
                </tr>
                <tr>
                    <td>Label Two:</td>
                    <td>
                        <input type="file" name="labeltwo" />
                    </td>
                </tr>
                <tr>
                    <td>Remark:</td>
                    <td>
                        <input type="text" name="remark" />
                    </td>
                </tr>
                <tr>
                    <td colspan="2">
                        <input type="submit" value="Add" class="btn btn-primary" />
                    </td>
                </tr>
            </table>

        </div>
        <div class="wrap-table">
            <table class="table table-striped table-bordered table-hover table-condensed">
                <tr>
                    <th>Index</th>
                    <th>Name</th>
                    <th>Price</th>
                    <th>Quantity</th>
                    <th>Amount</th>
                    <th>LabelOne</th>
                    <th>LabelTwo</th>
                    <th>Remark</th>
                </tr>
                {{range .Details}}
                <tr>
                    <td>{{.ID}}</td>
                    <td>{{.Name}}</td>
                    <td>{{.Price}}</td>
                    <td>{{.Quantity}}</td>
                    <td>{{getamount .Price .Quantity}}</td>
                    <td>{{if .LabelOne}}<img class="smallimg" src="data:image/jpg;base64,{{.LabelOne}}">{{else}}None{{end}}</td>
                    <td>{{if .LabelTwo}}<img class="smallimg" src="data:image/jpg;base64,{{.LabelTwo}}">{{else}}None{{end}}</td>
                    <td>{{.Remark}}</td>
                </tr>
                {{end}}
            </table>
        </div>
    </form>
    <script src="/static/js/jquery-1.12.3.min.js"></script>
    <script src="/static/jquery-ui-1.11.4/jquery-ui.min.js"></script>
    <script src="/static/bootstrap-3.3.6/js/bootstrap.min.js"></script>

    <script src="/static/js/common.js"></script>
</body>

</html>
