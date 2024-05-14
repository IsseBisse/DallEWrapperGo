import { useEffect, useState } from 'react'
import Box from '@mui/material/Box';
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import { CardActionArea } from '@mui/material';


function ImageThumbnail({ id, backendURL, setImageFocused }) {
    const [description, setDescription] = useState("")
    const [lowResImage, setLowResImage] = useState("")
    const [image, setImage] = useState("")

    useEffect(() => {
        if (lowResImage === "") {
            fetch(backendURL + "images/" + id)
            .then(response => response.json())
            .then(jsonData => {
                setDescription(jsonData.prompt)
                setLowResImage(jsonData.data)
            })
        }
      }, [])
    
    const card = (
        <>
        <CardActionArea onClick={() => {
            // fetch(backendURL + "images/" + id + "?isHighResolution")
            // .then(response => response.json())
            // .then(jsonData => setImageFocused(jsonData.data, description))
            setImageFocused(id, lowResImage, description)
        }}>
            <CardContent>
                <img 
                    src={lowResImage === "" ? "./loading.gif" : `data:image/jpeg;base64,${lowResImage}`}
                    width="100%"
                    
                    />
            </CardContent>
        </CardActionArea>
        </>
    );

    return (
        <Box 
            width="500" 
            height="500" 
        >
            <Card variant="outlined">{card}</Card>
        </Box>
    );
}

export default ImageThumbnail