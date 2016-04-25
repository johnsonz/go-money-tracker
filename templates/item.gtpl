{{template "header" .}} {{template "nav" .}}
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
                    <input type="text" name="purchaseddate" class="datepicker" />
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
                <th>Created Time</th>
                <th>Created By</th>
                <th colspan="4">Action</th>
            </tr>
            {{if .Items}}{{range .Items}}
            <tr>
                <td><span name="itemid">{{.ID}}</span></td>
                <td><span name="itemcate">{{.Subcategory.Category.Name}}</span></td>
                <td><span name="itemsubcate">{{.Subcategory.Name}}</span></td>
                <td><span name="itemstore">{{.Store}}</span></td>
                <td><span name="itemaddr">{{.Address}}</span></td>
                <td><span name="itempurdate">{{.PurchasedDate}}</span></td>
                <td><span name="itemamount">{{.Amount}}</span></td>
                <td><span name="itemreceipt">{{if .Receipt}}<img class="smallimg" src="data:image/jpg;base64,{{.Receipt}}">{{else}}None{{end}}</span></td>
                <td><span name="itemremark">{{.Remark}}</span></td>
                <td><span name="itemctime">{{.Operation.CreatedTime}}</span></td>
                <td><span name="itemcby">{{.Operation.CreatedBy}}</span></td>
                <td><a href="/item?id={{.ID}}" class="btn btn-link">View</td>
                <td><a href="javascript:" data-toggle="modal" data-target="#Modal-Item" class="btn btn-link edit">Edit</td>
                <td><a href="/detail?id={{.ID}}" class="btn btn-link">Add</td>
                <td><a href="/item?id={{.ID}}&action=del&page={{$.Pagination.Index}}" class="btn btn-link">Delete</td>
            </tr>
            {{end}}
            {{else}}
            <tr>
                <td colspan="13">
                    <blockquote>
                    <p>No data found.</p>
                    </blockquote>
                </td>
            </tr>
            {{end}}
        </table>
        <nav>
            <ul class="pagination">
                {{if le .Pagination.Previous 0}}
                <li class="disabled">
                    <a href="" aria-label="Previous">
                        <span aria-hidden="true">&laquo;</span>
                    </a>
                </li>
                {{else}}
                <li>
                    <a href="javascript:" onclick="this.href='/item?cid='+$('#category').val()+'&sid='+$('#subcategory').val()+'&page={{.Pagination.Previous}}'" aria-label="Previous">
                        <span aria-hidden="true">&laquo;</span>
                    </a>
                </li>
                {{end}} {{if gt (minus .Pagination.Index 2) 0}}
                <li><a href="javascript:" onclick="this.href='/item?cid='+$('#category').val()+'&sid='+$('#subcategory').val()+'&page={{minus .Pagination.Index 1}}'">{{minus .Pagination.Index 2}}</a></li>
                {{end}} {{if gt (minus .Pagination.Index 1) 0}}
                <!-- <li><a href="/item?page={{minus .Pagination.Index 1}}">{{minus .Pagination.Index 1}}</a></li> -->
                <li><a href="javascript:" onclick="this.href='/item?cid='+$('#category').val()+'&sid='+$('#subcategory').val()+'&page={{minus .Pagination.Index 1}}'">{{minus .Pagination.Index 1}}</a></li>
                {{end}}
                <li class="active"><a href="javascript:" onclick="this.href='/item?cid='+$('#category').val()+'&sid='+$('#subcategory').val()+'&page={{.Pagination.Index}}'">{{.Pagination.Index}}</a></li>
                {{if le (plus .Pagination.Index 1) .Pagination.Count}}
                <li><a href="javascript:" onclick="this.href='/item?cid='+$('#category').val()+'&sid='+$('#subcategory').val()+'&page={{plus .Pagination.Index 1}}'">{{plus .Pagination.Index 1}}</a></li>
                {{end}} {{if le (plus .Pagination.Index 2) .Pagination.Count}}
                <li><a href="javascript:" onclick="this.href='/item?cid='+$('#category').val()+'&sid='+$('#subcategory').val()+'&page={{plus .Pagination.Index 2}}'">{{plus .Pagination.Index 2}}</a></li>
                {{end}} {{if le .Pagination.Next .Pagination.Count}}
                <li>
                    <a href="javascript:" onclick="this.href='/item?cid='+$('#category').val()+'&sid='+$('#subcategory').val()+'&page={{.Pagination.Next}}'" aria-label="Next">
                        <span aria-hidden="true">&raquo;</span>
                    </a>
                </li>
                {{else}}
                <li class="disabled">
                    <a href="" aria-label="Next">
                        <span aria-hidden="true">&raquo;</span>
                    </a>
                </li>
                {{end}}

            </ul>
        </nav>
    </div>
    <input type="text" name="updatedid" id="updatedid"  hidden="hidden">
    <div class="modal fade" id="Modal-Item" tabindex="-1" role="dialog" aria-labelledby="ModalLabel">
        <div class="modal-dialog" role="document">
            <div class="modal-content">
                <div class="modal-header">
                    <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
                    <h4 class="modal-title" id="ModalLabel">Edit</h4>
                </div>
                <div class="modal-body">
                    <form>
                        <div class="form-group">
                            <label for="updatedcategory" class="control-label">Category:</label>
                            <select name="updatedcategory" class="form-control" id="updatedcategory">
                                {{range .Categories}}
                                <option value="{{.ID}}" {{if .Selected}}selected="selected" {{end}}>{{.Name}}</option>
                                {{end}}
                            </select>
                        </div>
                        <div class="form-group">
                            <label for="updatedsubcategory" class="control-label">Subcategory:</label>
                            <select name="updatedsubcategory" class="form-control" id="updatedsubcategory">
                            </select>
                        </div>
                        <div class="form-group">
                            <label for="store" class="control-label">Store:</label>
                            <input type="text" class="form-control" name="updatedstore" id="updatedstore">
                        </div>
                        <div class="form-group">
                            <label for="address" class="control-label">Address:</label>
                            <input type="text" class="form-control" name="updatedaddress" id="updatedaddress">
                        </div>
                        <div class="form-group">
                            <label for="updatedpurchaseddate" class="control-label">PurchasedDate:</label>
                            <input type="text" name="updatedpurchaseddate" id="updatedpurchaseddate" class="datepicker form-control" />
                        </div>
                        <div class="form-group">
                            <label for="receipt" class="control-label">Receipt:</label>
                            <div id="wrapreceip">

                            </div>
                            <input type="file" class="form-control" name="updatedreceipt" id="purchaseddatereceipt">
                        </div>
                        <div class="form-group">
                            <label for="remark" class="control-label">Remark:</label>
                            <input type="text" class="form-control" name="purchaseddateremark" id="purchaseddateremark">
                        </div>
                        <div class="form-group">
                            <label for="createdtime" class="control-label">Created Time:</label>
                            <input type="text" class="form-control" id="createdtime" disabled="disabled">
                        </div>
                        <div class="form-group">
                            <label for="createdby" class="control-label">Created By:</label>
                            <input type="text" class="form-control" id="createdby" disabled="disabled">
                        </div>
                    </form>
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-default" data-dismiss="modal">Cancel</button>
                    <input type="submit" name="update" value="Update" class="btn btn-primary" />
                </div>
            </div>
        </div>
    </div>
</form>
{{template "footer" .}}
