import { Flex, Spinner, Stack, Text, useColorModeValue } from "@chakra-ui/react";
import Item from "./Item";

const List = ({ tracklist, setTracklist }) => {
	const updateTrack = (updatedTrack) => {
		setTracklist((prevTracks) =>
			prevTracks.map((track) =>
				track.body === updatedTrack.body && track.artist === updatedTrack.artist
					? updatedTrack
					: track
			)
		);
	};

	return (

    <>
      {tracklist.length === 0 ? (
        <Stack alignItems={"center"} gap='3'>
          <Text fontSize={"xl"} textAlign={"center"} color={"gray.500"}>
            No results found
          </Text>
        </Stack>
      ) : (
        <Stack
          border={"1px"}
          borderRadius={"8px"}
          py={"4px"}
          px={"8px"}
          borderColor={useColorModeValue("gray.200", "gray.600")}
          gap={2}
        >
          {tracklist.map((track, i) => (
            <Item key={i} track={track} updateTrack={updateTrack}/>
          ))}
        </Stack>
      )}
    </>
	);
};

export default List;