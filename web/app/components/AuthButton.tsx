'use client';
import { signIn, signOut } from 'next-auth/react';
import React from 'react';

// ログインボタン
export const Login = () => {
  return (
    <button className='btn btn-blue' data-testid='login' onClick={() => signIn('cognito')}>
      サインイン
    </button>
  );
};

// ログアウトボタン
export const Logout = () => {
  return (
    <button className='btn btn-blue' onClick={() => signOut()}>
      サインアウト
    </button>
  );
};
