{{template "header" .}}
{{template "nav" .}}
    <form action="/item" method="POST" enctype="multipart/form-data">
        <div class="wrap-field">
            <table class="table table-condensed">
                <tr>
                    <td>Category:</td>
                    <td>
                        <select name="category" id="category">
                            {{range .Categories}}
                            <option value="{{.ID}}" {{if .Selected}}selected="selected" {{end}}>{{.Name}}</option>
                            {{end}}
                        </select>
                    </td>
                </tr>
                <tr>
                    <td>Subcategory:</td>
                    <td>
                        <select name="subcategory" id="subcategory">
                            {{range .Subcategories}}
                            <option value="{{.ID}}" {{if .Selected}}selected="selected" {{end}}>{{.Name}}</option>
                            {{end}}
                        </select>
                    </td>
                </tr>
                <tr>
                    <td>PurchasedDate:</td>
                    <td>
                        <input type="text" name="purchaseddate" id="datepicker" />
                    </td>
                </tr>
                <tr>
                    <td>Store:</td>
                    <td>
                        <input type="text" name="store" />
                    </td>
                </tr>
                <tr>
                    <td>Address:</td>
                    <td>
                        <input type="text" name="address" />
                    </td>
                </tr>
                <tr>
                    <td>Remark:</td>
                    <td>
                        <input type="text" name="remark" />
                    </td>
                </tr>
                <tr>
                    <td>Receipt Image:</td>
                    <td>
                        <input type="file" name="receiptimage" />
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
                    <th>Category</th>
                    <th>Subcategory</th>
                    <th>Store</th>
                    <th>Address</th>
                    <th>PurchasedDate</th>
                    <th>Amount</th>
                    <th>Receipt</th>
                    <th>Remark</th>
                </tr>
                {{range .Items}}
                <tr>
                    <td>{{.ID}}</td>
                    <td>{{.Subcategory.Category.Name}}</td>
                    <td>{{.Subcategory.Name}}</td>
                    <td>{{.Store}}</td>
                    <td>{{.Address}}</td>
                    <td>{{.PurchasedDate}}</td>
                    <td>{{.Amount}}</td>
                    <td>{{if .Receipt}}<img class="smallimg" src="data:image/jpg;base64,{{.Receipt}}">{{else}}None{{end}}</td>
                    <td>{{.Remark}}</td>
                </tr>
                {{end}}
            </table>
        </div>
    </form>
{{template "footer" .}}
