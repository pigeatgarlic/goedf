using System;
using System.Collections.Generic;
using System.Linq;
using System.Security.Claims;

namespace Authenticator.Models.Auth
{
    public class JwtOptions
    {
        public string Key { get; set; }

        public int ExpireHours {get;set;}
    }
}