:global SpeedTest do={
    :if (!any$url && [:typeof $url] != "str") do={
        :return "can't use that url bro!"
    }
    :local address $url;
    :local id [:rndnum from=10000000 to=99999999];
    :local cout ({});
    :for i from=0 to=4 do={
        :do {
            :local ping ([/tool ping google.com count=1 as-value]->"time");
            :set ($cout->$i) ([([:parse ([/tool fetch url="$address?seq=$i&id=$id&ping=$ping" mode=http http-method=get output=user as-value]->"data")])]); 
        } on-error={:put "err"}
    }
    :return ($cout->([:len $cout]-1));
}