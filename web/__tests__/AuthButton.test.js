import { render, screen, fireEvent } from '@testing-library/react';
import { Login, Logout } from '../app/components/AuthButton';
import '@testing-library/jest-dom';

describe('test for AuthButton', () => {
  it('Login', () => {
    render(<Login />);
    const LoginButton = screen.getByRole('button');

    expect(LoginButton).toBeInTheDocument();
    expect(LoginButton).toHaveTextContent('サインイン');
  });

  it('Logout', () => {
    render(<Logout />);
    const LogoutButton = screen.getByRole('button');
    expect(LogoutButton).toBeInTheDocument();
    expect(LogoutButton).toHaveTextContent('サインアウト');
  });
});
