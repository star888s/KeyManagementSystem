'use client';
import { signIn, signOut } from 'next-auth/react';
import React from 'react';
import { useRouter } from 'next/navigation';

// ログインボタン
export const Login = () => {
  return <button onClick={() => signIn('cognito')}>サインイン</button>;
};

// ログアウトボタン
export const Logout = () => {
  const router = useRouter();

  return <button onClick={() => signOut()}>サインアウト</button>;
};
