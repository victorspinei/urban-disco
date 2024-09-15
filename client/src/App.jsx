import { Container, Stack } from "@chakra-ui/react";
import { useState } from "react";
import Navbar from "./components/Navbar";
import MyForm from "./components/MyForm";
import List from "./components/List";

export const BASE_URL = import.meta.env.MODE === "development" ?
  "http://localhost:5000/api" : "/api";

function App() {
  const [tracklist, setTracklist] = useState([]);  // State to hold tracklist

  // Handle form submission
  const handleFormSubmit = async (e, query) => {
    e.preventDefault();

    try {
      const response = await fetch(`${BASE_URL}/tracklist/${encodeURIComponent(query)}`);
      if (!response.ok) {
        throw new Error("Failed to fetch tracklist");
      }

      const data = await response.json(); // Assuming the API returns JSON
      setTracklist(data); // Update the tracklist state with the fetched data

    } catch (error) {
      console.error("Error fetching tracklist:", error);
    }
  };

  return (
    <>
      <Stack h='100vh'>
        <Navbar />
        <Container>
          <MyForm handleFormSubmit={handleFormSubmit} />
          <List tracklist={tracklist} setTracklist={setTracklist}/>
        </Container>
      </Stack>
    </>
  );
}

export default App;
