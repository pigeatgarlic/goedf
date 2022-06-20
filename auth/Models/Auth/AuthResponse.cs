using System;
using System.Collections.Generic;
using System.Linq;
using Microsoft.AspNetCore.Identity;

namespace Authenticator.Models.Auth
{
    public class LoginResponse
    {
        public string? UserName { get; set; }
        public List<IdentityError>? Errors {get;set;}
        public string? Token { get; set; }
        public DateTime? ValidUntil { get; set; }


        public static LoginResponse GenerateSuccessful(string? username, string? token, DateTime? expiry)
        {
            return new LoginResponse()
            {
                Errors = null, 
                UserName = username,
                Token = token,
                ValidUntil = expiry,
            };
        }

        public static LoginResponse GenerateFailure(string username, IEnumerable<IdentityError> errcode)
        {
            return new LoginResponse()
            {
                Errors = errcode.ToList(),
                UserName = username
            };
        }
        public static LoginResponse GenerateFailure(string username, IdentityError errcode)
        {
            var err = new List<IdentityError>();
            err.Add(errcode);
            return new LoginResponse()
            {
                Errors = err,
                UserName = username
            };
        }
    }
}
