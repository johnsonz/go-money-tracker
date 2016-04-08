<html>
<head>
    <title></title>
    <link rel="stylesheet" href="/static/jquery-ui-1.11.4/jquery-ui.min.css">
    <script src="/static/js/jquery-1.12.3.min.js"></script>
    <script src="/static/jquery-ui-1.11.4/jquery-ui.min.js"></script>
    <script src="/static/js/common.js"></script>
</head>
<body>
    <form action="/item" method="POST">
        <div>
            <table>
                <tr>
                    <td>PurchasedDate:</td>
                    <td><input type="text" name="purchaseddate" id="datepicker"/></td>
                </tr>
                <tr>
                    <td>Store:</td>
                    <td><input type="text" name="store"/></td>
                </tr>
                <tr>
                    <td>Address:</td>
                    <td><input type="text" name="address"/></td>
                </tr>
                <tr>
                    <td>Remark:</td>
                    <td><input type="number" name="remark"/></td>
                </tr>
                <tr>
                    <td>Receipt Image:</td>
                    <td><input  type="file" name="receiptimage"/></td>
                </tr>
            </table>
            <input type="submit" value="Add"/>
        </div>
        <div>
            <table>
                <tr>
                    <th>Index</th>
                    <th>Category</th>
                    <th>Subcategory</th>
                    <th>Store</th>
                    <th>Address</th>
                    <th>PurchasedDate</th>
                    <th>Receipt</th>
                    <th>Remark</th>
                </tr>
                {{range .Items}}
                <tr>
                    <td>{{.ID}}</td>
                <!--    <td>{{.Category.Name}}</td>
                    <td>{{.Subcategory.Name}}</td>-->
                    <td>{{.Store}}</td>
                    <td>{{.Address}}</td>
                    <td>{{.PurchasedDate}}</td>
                    <td>{{.Receipt}}</td>
                    <td>{{.Remark}}</td>
                </tr>
                {{end}}
            </table>
        </div>
    </form>
</body>
</html>
