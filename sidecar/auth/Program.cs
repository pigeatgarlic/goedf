using Microsoft.AspNetCore.Hosting;
using Microsoft.Extensions.Options;
using Microsoft.Extensions.Hosting;
using Authenticator.Database.Context;
using Microsoft.AspNetCore.Identity;
using Authenticator.Models;
using Microsoft.Extensions.DependencyInjection;
using Authenticator.Models.User;
using Authenticator.Database.Seeding;
using System;


namespace Authenticator
{
    public class Program
    {
        public static void Main(string[] args)
        {
            IHost host = null;
            try
            {
                host = CreateHostBuilder(args) .Build();
                SeedDatabase(host);
            }
            catch(Exception ex)
            {
                Console.WriteLine($"{ex.Message} \n {ex.StackTrace}");
                return;
            }
            host.Run();
        }

        public static IHostBuilder CreateHostBuilder(string[] args) =>
            Host.CreateDefaultBuilder(args)
                .ConfigureWebHostDefaults(webBuilder =>
                {
                    webBuilder.UseStartup<Startup>();
                });

        static void SeedDatabase(IHost host)
        {
            using (var scope = host.Services.CreateScope())
            {
                var services = scope.ServiceProvider;

                var db = services.GetRequiredService<GlobalDbContext>();
                var userManager = services.GetRequiredService<UserManager<UserAccount>>();
                var roleManager = services.GetRequiredService<RoleManager<IdentityRole<int>>>();
                var systemconfig = services.GetRequiredService<IOptions<SystemConfig>>();

                
                AccountSeeder.SeedRoles(roleManager);
                AccountSeeder.SeedUserRole(userManager);
                AccountSeeder.SeedAdminUsers(userManager,db,systemconfig.Value);
            }
        }
    }
}
