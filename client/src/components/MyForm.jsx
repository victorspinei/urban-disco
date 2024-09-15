import { Button, Flex, Input, Spinner } from "@chakra-ui/react";
import { useState } from "react";
import { CiSearch } from "react-icons/ci";

const MyForm = ({ handleFormSubmit }) => {
  const [query, setQuery] = useState("");
  const [isPending, setIsPending] = useState(false);

  // Handle form submission
  const submitForm = async (e) => {
    e.preventDefault(); // Prevent the default form submission behavior
    setIsPending(true); // Show spinner
    await handleFormSubmit(e, query);  // Pass the query to the parent handler
    setIsPending(false); // Hide spinner after request is complete
  };

  return (
    <form onSubmit={submitForm}>
      <Flex gap={2} mb={8}>
        <Input
          type="text"
          placeholder="Type song name here..."
          value={query}
          onChange={(e) => setQuery(e.target.value)} // Update query state as the user types
        />
        <Button
          mx={2}
          type="submit"
          _active={{
            transform: "scale(.97)",
          }}
        >
          {isPending ? <Spinner size={"xs"} /> : <CiSearch size={30} />}
        </Button>
      </Flex>
    </form>
  );
};

export default MyForm;
