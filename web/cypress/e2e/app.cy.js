describe('Navigation', () => {
  it('should navigate to the about page', () => {
    // Start from the index page
    cy.visit('/');

    cy.get('h1').contains('KeyManagementSystem');

    cy.get('.btn').click();

    const id = Cypress.env('id');
    const pass = Cypress.env('pass');

    if (cy.contains('Sign in with Cognito')) {
      cy.contains('Sign in with Cognito').click();
    }

    cy.origin(
      Cypress.env('cognito_domain'),
      {
        args: {
          id,
          pass,
        },
      },
      ({ id, pass }) => {
        cy.get('input[name="username"]:visible').type(id);
        cy.get('input[name="password"]:visible').type(pass);
        cy.get('input[name="signInSubmitButton"]:visible').click();
      },
    );
    cy.wait(2000);

    if (cy.contains('Sign in with Cognito')) {
      cy.contains('Sign in with Cognito').click();
    }

    cy.get('main').contains('Hello!');
  });
});
