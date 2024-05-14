import { useState } from 'react'
import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import Grid from '@mui/material/Grid';
import MenuItem from '@mui/material/MenuItem';
import Modal from '@mui/material/Modal';
import Select from '@mui/material/Select';
import Typography from '@mui/material/Typography';

import GenerateModalRow from './GenerateModalRow';

const style = {
    position: 'absolute',
    top: '50%',
    left: '50%',
    transform: 'translate(-50%, -50%)',
    width: 600,
    bgcolor: 'background.paper',
    color: 'black',
    border: '2px solid #000',
    boxShadow: 24,
    p: 4,
  };

function GenerateModal({ open, setOpen, generateImage, setStylePrompt, setStyleUrl, setScenePrompt, setSceneUrl }) {
    const [numImages, setNumImages] = useState(1)
    const [size, setSize] = useState("1024x1024")
    
    return (
        <>
        <Modal
            open={open}
            onClose={() => setOpen(false)}
            aria-labelledby="modal-modal-title"
            aria-describedby="modal-modal-description"
        >
            <Box sx={style}>
                <Grid container spacing={2}>
                    <Grid item xs={6}>
                        Image size &nbsp;
                        <Select
                            id="size"
                            value={size}
                            label="Size"
                            onChange={(event) => setSize(event.target.value)}
                        >
                            <MenuItem value={"1024x1024"}>1024x1024</MenuItem>
                            <MenuItem value={"1792x1024"}>1792x1024</MenuItem>
                            <MenuItem value={"1024x1792"}>1024x1792</MenuItem>
                        </Select>
                    </Grid>
                    <Grid item xs={6}>
                        Number of images &nbsp;
                        <Select
                            id="num-images"
                            value={numImages}
                            label="Number of images"
                            onChange={(event) => {
                                setNumImages(event.target.value)
                            }}
                        >
                            <MenuItem value={1}>1</MenuItem>
                            <MenuItem value={2}>2</MenuItem>
                            <MenuItem value={3}>3</MenuItem>
                        </Select>
                    </Grid>
                </Grid>
                <Grid container spacing={2}>
                    <Grid item xs={6}>
                        <GenerateModalRow 
                            heading="Style"
                            description="The overall tone and artistic style of the image. Happy or sad. Dark or light. Painted or animated."
                            setPrompt={setStylePrompt}
                            setUrl={setStyleUrl}
                        />
                    </Grid>
                    <Grid item xs={6}>
                        <GenerateModalRow 
                            heading="Scene"
                            description="The objects and environments depicted in the image. People or animals. In a forrest or on a beach."
                            setPrompt={setScenePrompt}
                            setUrl={setSceneUrl}
                        />
                    </Grid>
                </Grid>
                <Button 
                    variant="contained"
                    size="large"
                    sx={{ margin: 3 }} 
                    onClick={() => generateImage(numImages, size)}
                >
                    Generate
                </Button>
            </Box>
        </Modal>
        </>
    );
}

export default GenerateModal