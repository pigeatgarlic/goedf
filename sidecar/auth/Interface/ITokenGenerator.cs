using Authenticator.Models.User;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Identity;
using System.Collections.Generic;
using System;

namespace Authenticator.Interfaces
{
    public interface ITokenGenerator
    {
        Task<string> GenerateUserJwt(UserAccount user);

        Task<UserAccount?> ValidateUserToken(string token);

        Task<List<string>> GetUserRoles(UserAccount account);

        DateTime GetTokenExpireTime();
    }
}
