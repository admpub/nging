function attachSearchNetworkDriverField(){
    $('#networkDriver').not('.tt-input').each(function(){
        searchNetworkDriverField(this);
    });
}
function searchNetworkDriverField(elem,size){
  if(size==null)size=500;
  $(elem).typeahead({
   hint: true, highlight: true, minLength: 0
  }, {source: function (query, sync, async) {
      $.ajax({
        url: BACKEND_URL+'/docker/base/index',
        type: 'get',
        data: {prop:'plugins.network'},
        dataType: 'json',
        async: false,
        success: function (r) {
          if(r.Code!=1) return App.message({text:r.Info,type:'error'});
          if(!r.Data) return;
          sync(r.Data.listData);
        }
      });
  },limit: size});
}
function attachSearchConfigFromField(){
    $('#configFromNetwork').not('.tt-input').each(function(){
        searchConfigFromField(this);
    });
}
function searchConfigFromField(elem,size){
  if(size==null)size=500;
  $(elem).typeahead({
   hint: true, highlight: true, minLength: 0
  }, {source: function (query, sync, async) {
      $.ajax({
        url: BACKEND_URL+'/docker/base/network/index',
        type: 'get',
        data: {op:'ajaxList'},
        dataType: 'json',
        async: false,
        success: function (r) {
          if(r.Code!=1) return App.message({text:r.Info,type:'error'});
          if(!r.Data) return;
          sync(r.Data.listData);
        }
      });
  },limit: size});
}
$(function(){
    attachSearchNetworkDriverField();
    attachSearchConfigFromField();
})