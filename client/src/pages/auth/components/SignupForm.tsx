import { useState } from "react";
import { useForm } from "react-hook-form";
import { useMutation } from "@tanstack/react-query";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";

import { Input } from "../../../components/ui/input";
import { Button } from "../../../components/ui/button";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "../../../components/ui/form";
import { useAuth } from "../../../hooks/useAuth.tsx";
import { SignupRequest, SignupResponse } from "../../../hooks/AuthContext.tsx";
import ApiError from "../../../api/error.ts";

const CompanyNameInfo =
  "Set the name of your company, this is the name all of your employees will see themselves under when they sign in.";
const FirstNameInfo = "The first name of the initial admin user.";
const LastNameInfo = "The last name of the initial admin user.";
const EmailInfo =
  "The email address used here will be considered a super admin account and will have permissions to create, read, update, and delete all content. Permissions can be tailored later to give other users different privilages.";

const formSchema = z.object({
  name: z.string().min(4, {
    message: "Your company name should be at least 4 characters long",
  }),
  firstName: z.string().min(1, "A first name must be provided"),
  lastName: z.string().min(1, "A last name must be provided"),
  email: z.string().email(),
  password: z.string().min(6, {
    message: "Your password must be at least 6 characters long",
  }),
});

type FormSchema = z.infer<typeof formSchema>;

interface SignupFormParams {
  onSignupComplete: () => void;
}

function SignupForm({ onSignupComplete }: SignupFormParams) {
  const { signup } = useAuth();
  const [error, setError] = useState<null | ApiError>(null);
  const [emailPlaceholder, setEmailPlaceholder] = useState(
    "your.name@company.com"
  );

  const { control, getValues, handleSubmit, ...form } = useForm<FormSchema>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: "Company",
      firstName: "Mike",
      lastName: "Murphy",
      email: "mikhl90@gmail.com",
      password: "password",
    },
  });

  const { mutateAsync: signupMutation, isPending: isSignupPending } =
    useMutation<SignupResponse, ApiError, SignupRequest>({
      mutationFn: signup,
      onSuccess: () => {
        console.log("onSuccess");
        onSignupComplete();
      },
      onError: (error) => {
        console.log("onError");
        console.error(error);
        setError(error);
      },
    });

  async function onSubmit(values: FormSchema) {
    await signupMutation(values).catch(() => {});
  }

  function onCompanyNameChange(companyName: string) {
    setEmailPlaceholder(companyName ? `your.name@${companyName.toLowerCase().replace(/\s+/g, '')}.com` : "your.name@company.com");
  }

  return (
    <>
      {error && (
        <section className="space-y-2 bg-red-200 shadow mb-8 p-4 rounded-lg text-red-950">
          <p className="font-semibold">
            There has been an error signing you up:
          </p>
          <p>{error.message}</p>
        </section>
      )}
      <Form { ...{ control, getValues, handleSubmit, ...form } }>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <FormField
            control={control}
            name="name"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Company name</FormLabel>
                <FormControl onChange={() => onCompanyNameChange(getValues("name"))}>
                  <Input
                    autoFocus={true}
                    placeholder="Your Company LTD"
                    data-testid="company-name"
                    {...field}
                  />
                </FormControl>
                <FormMessage />
                <FormDescription className="hidden">
                  {CompanyNameInfo}
                </FormDescription>
              </FormItem>
            )}
          />
          <FormField
            control={control}
            name="firstName"
            render={({ field }) => (
              <FormItem>
                <FormLabel>First name</FormLabel>
                <FormControl>
                  <Input type="text" placeholder="Your first name" {...field} />
                </FormControl>
                <FormMessage />
                <FormDescription className="hidden">
                  {FirstNameInfo}
                </FormDescription>
              </FormItem>
            )}
          />
          <FormField
            control={control}
            name="lastName"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Last name</FormLabel>
                <FormControl>
                  <Input type="text" placeholder="Your last name" {...field} />
                </FormControl>
                <FormMessage />
                <FormDescription className="hidden">
                  {LastNameInfo}
                </FormDescription>
              </FormItem>
            )}
          />
          <FormField
            control={control}
            name="email"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Email</FormLabel>
                <FormControl>
                  <Input
                    type="email"
                    placeholder={emailPlaceholder}
                    data-testid="user-email"
                    {...field}
                  />
                </FormControl>
                <FormMessage />
                <FormDescription className="hidden">
                  {EmailInfo}
                </FormDescription>
              </FormItem>
            )}
          />
          <FormField
            control={control}
            name="password"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Password</FormLabel>
                <FormControl>
                  <Input
                    type="password"
                    placeholder="* * * * * * * * * * * * * *"
                    {...field}
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <Button type="submit" disabled={isSignupPending} loading={isSignupPending}>
            Sign up
          </Button>
        </form>
      </Form>
    </>
  );
}

export default SignupForm;
