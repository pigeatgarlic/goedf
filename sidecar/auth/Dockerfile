FROM mcr.microsoft.com/dotnet/sdk:6.0 AS base

WORKDIR /src
COPY . .

ENV PG_USER=postgres
ENV PG_PASS=CVxL7CiKgh
ENV PG_SERVER=postgresql-ha-1653545151-pgpool.postgres.svc.cluster.local


EXPOSE 80

RUN dotnet tool install --global dotnet-ef

RUN dotnet restore
RUN dotnet tool restore

RUN dotnet build

RUN chmod +x /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]