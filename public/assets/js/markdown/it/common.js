function markdownParse(box,isContainer){
    if(typeof(window.markdownit)=='undefined')return;
    if(isContainer!=false) box=box.find('.markdown-code');
    var md=markdownItInstance();
    box.each(function(){
        $(this).html(md.render($.trim($(this).html())));
        $(this).find("pre > code").each(function(){
            $(this).parent("pre").addClass("prettyprint linenums");
        });
        if(typeof(prettyPrint)!=="undefined") prettyPrint();
    });
}
function markdownItInstance(){
    var md = window.markdownit();
    if(typeof(window.markdownitEmoji)!='undefined') md.use(window.markdownitEmoji);
    return md;
}