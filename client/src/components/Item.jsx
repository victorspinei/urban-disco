import { Badge, Box, Flex, Text, useColorModeValue, Spinner } from "@chakra-ui/react";
import { IoMdDownload } from "react-icons/io";
import { IoMusicalNotesOutline } from "react-icons/io5";
import { IoMusicalNotesSharp } from "react-icons/io5";
import { FaCheckCircle } from "react-icons/fa";
import { BASE_URL } from "../App";
import { useState } from "react";

const Item = ({ track, updateTrack }) => {
  const [isDownloading, setIsDownloading] = useState(false);

  const handleDownload = async (track) => {
    setIsDownloading(true);
    try {
      const response = await fetch(`${BASE_URL}/song?name=${encodeURIComponent(track.body)}&artist=${encodeURIComponent(track.artist)}`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      });

      if (!response.ok) {
        throw new Error('Download failed');
      }

      // Handle the file download (e.g., create a link and click it)
      const blob = await response.blob();
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `${track.body}.mp3`; // Adjust filename as needed
      document.body.appendChild(a);
      a.click();
      a.remove();

      // Update the track's downloaded status
      updateTrack({ ...track, downloaded: true });
    } catch (error) {
      console.error('Error downloading track:', error);
    } finally {
      setIsDownloading(false);
    }
  };

  return (
    <Flex gap={2} alignItems={"center"}>
      <Flex
        flex={1}
        alignItems={"center"}
        borderColor={"gray.600"}
        p={2}
        borderRadius={"lg"}
        justifyContent={"space-between"}
      >
        <Flex gap={2}>
          <Flex justifyContent={"center"} alignItems={"center"}>
            {useColorModeValue(<IoMusicalNotesOutline />, <IoMusicalNotesSharp />)}
          </Flex>
          <Text color={useColorModeValue("black", "white")}>{track.body}</Text>
        </Flex>
        {track.downloaded ? (
          <Badge ml='1' colorScheme='green'>Downloaded</Badge>
        ) : (
          <Badge ml='1' colorScheme='yellow'>Not Downloaded</Badge>
        )}
      </Flex>
      {isDownloading ? (
        <Box paddingRight={"4px"} color={"yellow.500"} cursor={"pointer"}>
          <Spinner size={"sm"} />
        </Box>
      ) : !track.downloaded ? (
        <Box paddingRight={"4px"} color={"yellow.500"} cursor={"pointer"} onClick={() => handleDownload(track)}>
          <IoMdDownload size={20} />
        </Box>
      ) : (
        <Box paddingRight={"4px"} color={"green.500"} cursor={"pointer"}>
          <FaCheckCircle size={20} />
        </Box>
      )}
    </Flex>
  );
};

export default Item;
