
function removeIPv4Domain(a,k) {
  var c = a.parent(), i = a.data('index');
  var lastIndex = c.find('.input-group[data-index]:last').data('index');
  a.remove();
  if (lastIndex == i) {
    return;
  }
  var prefix = 'DNSServices['+k+'][IPv4Domains][';
  resortDomain(c, prefix);
}
function removeIPv6Domain(a,k) {
  var c = a.parent(), i = a.data('index');
  var lastIndex = c.find('.input-group[data-index]:last').data('index');
  a.remove();
  if (lastIndex == i) {
    return;
  }
  var prefix = 'DNSServices['+k+'][IPv6Domains][';
  resortDomain(c, prefix);
}
function resortDomain(c, prefix){
    c.find('.input-group[data-index]').each(function(index){
        var $this = $(this);
        var oldIndex = $this.data('index');
        if(index == oldIndex) return;
        $this.attr('data-index',index);
        $this.find('[name^="'+prefix+'"]').each(function(){
          var name = $(this).attr('name');
          name = name.substr(0, prefix.length);
          name = name.replace(/^[0-9]+\]/,index+']');
          name = prefix+name;
          $(this).attr('name',name);
        })
    });
}
function addIPv4Domain(a,k) {
  var t = $('#tmpl-domain-row').html();
  var c = a.parent();
  var i = c.find('.input-group[data-index]:last').data('index')+1;
  t = t.replace(/\{=domainK=\}/g,i);
  t = t.replace(/\{=k=\}/g,k);
  t = t.replace(/IPv6Domains/g,'IPv4Domains');
  t = t.replace(/removeIPv6Domain/g,'removeIPv4Domain');
  c.append(t);
}
function addIPv6Domain(a,k) {
  var t = $('#tmpl-domain-row').html();
  var c = a.parent();
  var i = c.find('.input-group[data-index]:last').data('index')+1;
  t = t.replace(/\{=domainK=\}/g,i);
  t = t.replace(/\{=k=\}/g,k);
  c.append(t);
}
$(function(){
    // $('#ddns-form').off().on('submit',function(e){
    //     e.preventDefault();
    //     $.post(window.location.href,$(this).serialize(),function(r){
    //         if(r.Code==1){
    //             $('#search-result').text(r.Data);
    //             return;
    //         }
    //     },'json');
    // });
    $('input[name="IPv6[Type]"],input[name="IPv4[Type]"]').on('click',function(){
      var name = $(this).attr('name');
      $('div[rel="'+this.id+'"]').show();
      $('input[name="'+name+'"]:not(:checked)').each(function(){
        $('div[rel="'+this.id+'"]').hide();
      });
    });
    $('input[name="IPv6[Type]"]:checked,input[name="IPv4[Type]"]:checked').trigger('click');
});