<html>
<head>
    <title></title>
</head>
<body>
    <form action="/item" method="POST">
        <div>
            <table>
                <tr>
                    <td>Start Date:</td>
                    <td><input type="text" name="startDate"/></td>
                    <td>End Date:</td>
                    <td><input type="text" name="endDate"/></td>
                </tr>
                <tr>
                    <td>Category:</td>
                    <td><select name="category">
                        {{range .Categories}}
                            <option value="{{.ID}}" {{if Selected}}selected="selected"{{end}}>{{.Name}}</option>
                        {{end}}
                    </select></td>
                    <td>Subcategory:</td>
                    <td><select name="subcategory">
                        {{range .Subcategories}}
                            <option value="{{.ID}}" {{if Selected}}selected="selected"{{end}}>{{.Name}}</option>
                        {{end}}
                    </select></td>
                </tr>
                <tr>
                    <td>Store:</td>
                    <td><select name="store">
                        {{range .Stores}}
                            <option value="{{.Name}}" {{if Selected}}selected="slected"{{end}}>{{.Name}}</option>
                        {{end}}
                    </select></td>
                    <td>Address:</td>
                    <td><input type="text" name="address"/></td>
                </tr>
                <tr>
                    <td>Min Amount:</td>
                    <td><input type="number" name="minamount"/></td>
                    <td>Max Amount:</td>
                    <td><input type="number" name="maxamount"/></td>
                </tr>
            </table>
        </div>
        <div>
            
        </div>
    </form>
</body>
</html>
