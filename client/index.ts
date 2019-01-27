import Vault from 'node-vault';

const vault = Vault({ endpoint: 'http://localhost:8200' });

const jwt = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJkZW1vIiwiZXhwIjoxNTQ4NjA5NTgwLCJpYXQiOjE1NDg2MDU5ODAsImlzcyI6InRlc3QiLCJzdWIiOiJtYXJjdXMiLCJhcHBfbWV0YWRhdGEiOnsicm9sZXMiOlsiYWRtaW5zIl19fQ.rf1Z8YZgPHz3R9jY2j0FYBMVxJTj3WWiThn5SCctTerDScDS8NmMfCOlZymFmkkNvfAT-7HPTm0nP1S7UVgjPHtnYoRX6G_ZXq1YboTEeXewsIpe1l-yp9HLO0w7BX3o-66wsOYPLOWxTqVCDjovlx5iRizklEFZSyMECFIa0RBjJ7vo0PRhvz8PhWA_7QZlzC4mRGyzhTafqotW8h3caO1K8DO1sYnTT0GW-4UBTulXUYEFHLE869EGoeFKotBxakv3dlHAfMS1vZrPsP0AK5AzmFUsPl_caFg3lHDnju8CVOA0-MbTrS25LwVCUbYP3x8i5AXzMzgdpNCLkMXrSg"

interface LoginResponse {
  auth: {
      client_token: string,
  }
}

interface ReadResponse<T = { [k: string]: any }> {
  data: T
}

type KVResponse = ReadResponse<{ data: { [k: string]: string } }>;

const login = async (): Promise<string> => {
  const res: LoginResponse = await vault.request({
    path: '/auth/jwt/login',
    method: 'POST',
    json: {
      jwt,
      role: "demo",
    },
  });
  return res.auth.client_token;
}

const main = async (path: string) => {
  const token = await login();
  vault.token = token;

  const pathTransformed = path.split('.').join('/');
  const data: KVResponse = await vault.read(`secret/data/${pathTransformed}`);
  console.log("val:", data.data.data.value);
}

const path = process.argv[2];
main(path).catch(console.error);
