import { fireEvent, render, screen } from "@testing-library/react";
import { describe, it, expect } from "vitest";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import SignupForm from "./SignupForm";
import {ReactNode} from "react";

const queryClient = new QueryClient()

const renderWithClient = (ui: ReactNode) => {
  return render(
    <QueryClientProvider client={queryClient}>
      {ui}
    </QueryClientProvider>
  );
};

describe("LoginForm component", () => {
  it("sets the placeholder to 'your.name@company.com' when the company name is empty", () => {
    const dummyFn = () => {};
    renderWithClient(<SignupForm onSignupComplete={dummyFn} />);

    const companyInput = screen.getByTestId("company-name") as HTMLInputElement;
    const emailInput = screen.getByTestId("user-email") as HTMLInputElement;

    // Has default value
    expect(emailInput).toHaveAttribute("placeholder", "your.name@company.com");

    // Changes when company name is edited
    fireEvent.change(companyInput, { target: { value: "My Company" } });
    expect(emailInput).toHaveAttribute("placeholder", "your.name@mycompany.com");

    // Goes back to the default when the company name is blank again
    fireEvent.change(companyInput, { target: { value: "" } });
    expect(emailInput).toHaveAttribute("placeholder", "your.name@company.com");

    // Handles multiple spaces
    fireEvent.change(companyInput, { target: { value: "Longer Company Name LTD" } })
    expect(emailInput).toHaveAttribute("placeholder", "your.name@longercompanynameltd.com");
  });
});
