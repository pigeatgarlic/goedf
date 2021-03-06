using System;
using System.ComponentModel.DataAnnotations;

namespace Authenticator.Models.Auth
{
    public class LoginModel
    {
        [Required]
        public string UserName { get; set; }

        [Required]
        public string Password { get; set; }
    }

    public class UpdatePasswordModel
    {
        public string? Old { get; set; }

        [Required]
        public string New { get; set; }
    }

    public class RegisterModel : LoginModel
    {
        [Required]
        [EmailAddress]
        public string Email { get; set; }

        [Required]
        public string FullName { get; set; }

        public DateTime? DateOfBirth { get; set; }

        public string? Jobs { get; set; }

        [Required]
        public string PhoneNumber { get; set; }
    }
}