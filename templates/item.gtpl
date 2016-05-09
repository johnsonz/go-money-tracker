{{template "header" .}} {{template "nav" .}}
<form action="/item" method="POST" enctype="multipart/form-data">
    <div class="wrap-primary">
        <button type="button" class="btn btn-primary" data-toggle="modal" data-target="#Modal-Create-Item">Create new item</button>
    </div>
    <div class="wrap-table-item">
        <table class="table table-striped table-bordered table-hover table-condensed">
            <tr>
                <th>#</th>
                <th>Category</th>
                <th>Subcategory</th>
                <th>Store</th>
                <th>Address</th>
                <th>P.Date</th>
                <th>Amount</th>
                <th>Receipt</th>
                <th>Remark</th>
                <th>C.Time</th>
                <th>C.By</th>
                <th colspan="4">Action</th>
            </tr>
            {{if .Items}}{{range $index,$item:= .Items}}
            <tr>
                <td hidden="hidden"><span name="itemid">{{.ID}}</span></td>
                <td><span name="itemIndex">{{plus $index 1}}</span></td>
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
                <td><a href="/detail?id={{.ID}}" class="btn btn-link">View</td>
                <td><a href="javascript:" data-toggle="modal" data-target="#Modal-Item" class="btn btn-link edit">Edit</td>
                <td><a href="/detail?id={{.ID}}" class="btn btn-link">Add</td>
                <!-- <td><a href="/item?id={{.ID}}&action=del&page={{$.Pagination.Index}}" class="btn btn-link">Delete</td> -->
                <td><a href="javascript:" class="btn btn-link" data-toggle="modal" data-target="#Modal-Confirm-Item">Delete</td>

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
                    <a href="/item?page={{.Pagination.Previous}}" aria-label="Previous">
                        <span aria-hidden="true">&laquo;</span>
                    </a>
                </li>
                {{end}} {{if gt (minus .Pagination.Index 2) 0}}
                <li><a href="/item?page={{minus .Pagination.Index 1}}">{{minus .Pagination.Index 2}}</a></li>
                {{end}} {{if gt (minus .Pagination.Index 1) 0}}
                <!-- <li><a href="/item?page={{minus .Pagination.Index 1}}">{{minus .Pagination.Index 1}}</a></li> -->
                <li><a href="/item?page={{minus .Pagination.Index 1}}">{{minus .Pagination.Index 1}}</a></li>
                {{end}}
                <li class="active"><a href="/item?page={{.Pagination.Index}}" id="pageIndex">{{.Pagination.Index}}</a></li>
                {{if le (plus .Pagination.Index 1) .Pagination.Count}}
                <li><a href="/item?page={{plus .Pagination.Index 1}}">{{plus .Pagination.Index 1}}</a></li>
                {{end}} {{if le (plus .Pagination.Index 2) .Pagination.Count}}
                <li><a href="/item?page={{plus .Pagination.Index 2}}">{{plus .Pagination.Index 2}}</a></li>
                {{end}} {{if le .Pagination.Next .Pagination.Count}}
                <li>
                    <a href="/item?page={{.Pagination.Next}}" aria-label="Next">
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
    <div class="modal fade" id="Modal-Confirm-Item" tabindex="-1" role="dialog" aria-labelledby="ModalLabel" aria-hidden="true">
        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-header">
                    Confirm
                </div>
                <div class="modal-body">
                    This record will be permanently deleted and cannot be recovered. Are you sure?
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-default" data-dismiss="modal">Cancel</button>
                    <a class="btn btn-danger btn-ok itemdel">Delete</a>
                </div>
            </div>
        </div>
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
                            <label for="updatedstore" class="control-label">Store:</label>
                            <input type="text" class="form-control" name="updatedstore" id="updatedstore">
                        </div>
                        <div class="form-group">
                            <label for="updatedaddress" class="control-label">Address:</label>
                            <input type="text" class="form-control" name="updatedaddress" id="updatedaddress">
                        </div>
                        <div class="form-group">
                            <label for="updatedpurchaseddate" class="control-label">PurchasedDate:</label>
                            <input type="text" name="updatedpurchaseddate" id="updatedpurchaseddate" class="datepicker form-control" />
                        </div>
                        <div class="form-group">
                            <label for="updatedreceipt" class="control-label">Receipt:</label>
                            <div id="wrapreceip"></div>
                            <input type="file" class="form-control" name="updatedreceipt" id="updatedreceipt">
                        </div>
                        <div class="form-group">
                            <label for="updatedremark" class="control-label">Remark:</label>
                            <input type="text" class="form-control" name="updatedremark" id="updatedremark">
                        </div>
                        <div class="form-group">
                            <label for="createdtime" class="control-label">Created Time:</label>
                            <input type="text" class="form-control" id="createdtime" disabled="disabled">
                        </div>
                        <div class="form-group">
                            <label for="createdby" class="control-label">Created By:</label>
                            <input type="text" class="form-control" id="createdby" disabled="disabled">
                        </div>
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-default" name="cancel" data-dismiss="modal">Cancel</button>
                    <input type="submit" name="update" value="Update" class="btn btn-primary" />
                </div>
            </div>
        </div>
    </div>
    <div class="modal fade" id="Modal-Create-Item" tabindex="-1" role="dialog" aria-labelledby="ModalLabel">
        <div class="modal-dialog" role="document">
            <div class="modal-content">
                <div class="modal-header">
                    <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
                    <h4 class="modal-title" id="ModalLabel">Create new item</h4>
                </div>
                <div class="modal-body">
                        <div class="form-group">
                            <label for="createdcategory" class="control-label">Category:</label>
                            <select name="createdcategory" class="form-control" id="createdcategory">
                                {{range .Categories}}
                                <option value="{{.ID}}" {{if .Selected}}selected="selected" {{end}}>{{.Name}}</option>
                                {{end}}
                            </select>
                        </div>
                        <div class="form-group">
                            <label for="createdsubcategory" class="control-label">Subcategory:</label>
                            <select name="createdsubcategory" class="form-control" id="createdsubcategory">
                                    {{range .Subcategories}}
                                    <option value="{{.ID}}" {{if .Selected}}selected="selected" {{end}}>{{.Name}}</option>
                                    {{end}}
                                </select>
                            </select>
                        </div>
                        <div class="form-group">
                            <label for="createdstore" class="control-label">Store:</label>
                            <input type="text" class="form-control" name="createdstore" id="createdstore">
                        </div>
                        <div class="form-group">
                            <label for="createdaddress" class="control-label">Address:</label>
                            <input type="text" class="form-control" name="createdaddress" id="createdaddress">
                        </div>
                        <div class="form-group">
                            <label for="createdpurchaseddate" class="control-label">PurchasedDate:</label>
                            <input type="text" name="createdpurchaseddate" id="createdpurchaseddate" class="datepicker form-control" />
                        </div>
                        <div class="form-group">
                            <label for="createdreceipt" class="control-label">Receipt:</label>
                            <div id="wrapreceip"></div>
                            <input type="file" class="form-control" name="createdreceipt" id="createdreceipt">
                        </div>
                        <div class="form-group">
                            <label for="createdremark" class="control-label">Remark:</label>
                            <input type="text" class="form-control" name="createdremark" id="createdremark">
                        </div>
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-default" name="cancel" data-dismiss="modal">Cancel</button>
                    <input type="submit" name="create" value="Create" class="btn btn-primary" />
                </div>
            </div>
        </div>
    </div>
</form>
{{template "footer" .}}
