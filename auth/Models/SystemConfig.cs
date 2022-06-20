using Authenticator.Models.Auth;
using System.Collections.Generic;

namespace Authenticator.Models
{
    public class SystemConfig
    {
        public string GoogleOauthID {get;set;}
        public LoginModel AdminLogin { get; set; }
    }
}
