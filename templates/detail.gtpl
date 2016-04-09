<html>
<head>
    <title></title>
    <script src="/static/js/jquery-1.12.3.min.js"></script>
    <script src="/static/js/common.js"></script>
</head>
<body>
    <form action="/detail" method="POST" enctype="multipart/form-data">
        <input type="hidden" name="itemid" value="{{.ItemID}}" readonly="readonly"/>
        <div>
            <table>
                <tr>
                    <td>Name:</td>
                    <td><input type="text" name="name"/></td>
                </tr>
                <tr>
                    <td>Price:</td>
                    <td><input type="number" step="0.01" min="0" name="price"/></td>
                </tr>
                <tr>
                    <td>Quantity:</td>
                    <td><input type="number" min="0" name="quantity" value="1"/></td>
                </tr>
                <tr>
                    <td>Label One:</td>
                    <td><input type="file" name="labelone"/></td>
                </tr>
                <tr>
                    <td>Label Two:</td>
                    <td><input type="file" name="labeltwo"/></td>
                </tr>
                <tr>
                    <td>Remark:</td>
                    <td><input type="text" name="remark"/></td>
                </tr>

            </table>
            <input type="submit" value="Add"/>
        </div>
        <div>
            <table>
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
                    <td><img src="data:image/jpg;base64,{{.LabelOne}}"></td>
                    <td><img src="data:image/jpg;base64,{{.LabelTwo}}"></td>
                    <td>{{.Remark}}</td>
                </tr>
                {{end}}
            </table>
        </div>
    </form>
</body>
</html>
