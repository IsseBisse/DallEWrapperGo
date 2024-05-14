import Box from '@mui/material/Box';
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import Typography from '@mui/material/Typography';
import TextField from '@mui/material/TextField';


function GenerateModalRow({ heading, description, setPrompt, setUrl }) {
    const card = (
        <>
        <CardContent>
            <Typography variant="h5" component="div">
                    {heading}
            </Typography>
            <Typography variant="body2">
                {description}
                <br/>
                <br/>
                Please provide either a description or a link to a similar image.
            </Typography>
            <TextField 
                id={heading + "-prompt"} 
                label="Description" 
                variant="standard" 
                multiline
                // minRows={4}
                maxRows={4}
                onChange={e => setPrompt(e.target.value)}
            />
            <Typography variant="body2">
            or
            </Typography>
            <TextField 
                id={heading + "-url"} 
                label="Link" 
                variant="standard" 
                onChange={e => setUrl(e.target.value)}
            />
        </CardContent>
        </>
    );

    return (
        <Box 
            width="200" 
            // height="500" 
        >
            <Card variant="outlined">{card}</Card>
        </Box>
    );
}

export default GenerateModalRow