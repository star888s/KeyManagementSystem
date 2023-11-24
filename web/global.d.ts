declare namespace NodeJS {
  // 環境変数名の定義
  interface ProcessEnv {
    readonly COGNITO_CLIENT_ID: string;
    readonly COGNITO_CLIENT_SECRET: string;
    readonly COGNITO_ISSUER: string;
  }
}
