using System.ComponentModel.DataAnnotations;
using System.Collections.Generic;
using System;

namespace Authenticator.Models.Auth
{
    public class ValidateRequest
    {
        [Required]
        public string? jwtToken {get;set;}

        public string? googleToken {get;set;}
    }

    public class ValidateResponse 
    {
        [Required]
        public string UserID {get;set;}

        [Required]
        public List<string> Roles {get;set;}

        [Required]
        public DateTime ValidatedAt{get;set;}
    }
}