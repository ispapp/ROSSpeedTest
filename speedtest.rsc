:global SpeedTest do={
    :if (!any$url && [:typeof $url] != "str") do={
        :return "can't use that url bro!"
    }
    :local address $url;
    :local id [:rndnum from=10000000 to=99999999];
    :local cout ({});
    :local data [:rndstr from="abcdef%^&" length=100];
    :for i from=0 to=4 do={
        :do {
            :set ($cout->$i) ([([:parse ([/tool fetch url="$address?seq=$i&id=$id" http-data=$data  mode=http http-method=post output=user as-value]->"data")])]); 
        } on-error={}
    }
    :return ($cout->([:len $cout]-1));
}