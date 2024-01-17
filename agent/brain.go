package agent

import (
	. "gorgonia.org/gorgonia"
	"gorgonia.org/tensor"
)

type Layer struct {
	W   *Node
	Act func(x *Node) (*Node, error)
}

func (l *Layer) fwd(x *Node) (*Node, error) {
	xw := Must(Mul(x, l.W))
	if l.Act == nil {
		return xw, nil
	}
	activation, err := l.Act(xw)

	return activation, err
}

type Brain struct {
	g *ExprGraph
	x *Node
	y *Node
	l []Layer

	pred    *Node
	predVal Value
}

func NewBrain(numNeurons int) *Brain {
	g := NewGraph()

	x := NewMatrix(g, of, WithShape(1, 11), WithName("X"), WithInit(Zeroes()))
	y := NewMatrix(g, of, WithShape(1, 4), WithName("Y"), WithInit(Zeroes()))
	l := []Layer{
		{W: NewMatrix(g, tensor.Float32, WithShape(11, numNeurons), WithName("L0W"), WithInit(GlorotU(1.0))), Act: Rectify},
		{W: NewMatrix(g, tensor.Float32, WithShape(numNeurons, 20), WithName("L1W"), WithInit(GlorotU(1.0))), Act: Rectify},
		{W: NewMatrix(g, tensor.Float32, WithShape(20, 50), WithName("L2W"), WithInit(GlorotU(1.0))), Act: Rectify},
		{W: NewMatrix(g, tensor.Float32, WithShape(50, 4), WithName("L3W"), WithInit(GlorotU(1.0)))},
	}
	return &Brain{
		g: g,
		x: x,
		y: y,
		l: l,
	}
}

func (nn *Brain) forward(x *Node) (*Node, error) {
	var err error
	pred := x
	pred, err = nn.l[0].fwd(pred)
	if err != nil {
		return nil, err
	}
	pred, err = nn.l[1].fwd(pred)
	if err != nil {
		return nil, err
	}
	pred, err = nn.l[2].fwd(pred)
	if err != nil {
		return nil, err
	}
	pred, err = SoftMax(pred)
	if err != nil {
		return nil, err
	}

	return pred, nil
}

func (nn *Brain) learnables() Nodes {
	retVal := make(Nodes, 0, len(nn.l))
	for _, l := range nn.l {
		retVal = append(retVal, l.W)
	}
	return retVal
}

func (nn *Brain) model() []ValueGrad { return NodesToValueGrads(nn.learnables()) }

func (nn *Brain) cons() (pred *Node, err error) {
	pred = nn.x
	for _, l := range nn.l {
		if pred, err = l.fwd(pred); err != nil {
			return nil, err
		}
	}
	nn.pred = pred
	Read(nn.pred, &nn.predVal)

	cost := Must(Mean(Must(Square(Must(Sub(nn.y, pred))))))
	if _, err = Grad(cost, nn.learnables()...); err != nil {
		return nil, err
	}

	return pred, nil
}

func (nn *Brain) Let2(xs [11]float32, y float32) {
	xval := nn.x.Value().Data().([]float32)
	yval := nn.y.Value().Data().([]float32)

	//  overwrite the data
	for i := range xval {
		xval[i] = xs[i]
	}
	for i := range yval {
		yval[i] = 0
	}

	// For now, assume all the magic ends up in the first value of this array
	// TODO gain more clarity on out input state vector x related to output vector y
	yval[0] = y
}

func (nn *Brain) Let1(x [11]float32) {
	xval := nn.x.Value().Data().([]float32)
	// overwrite the data
	for i := range xval {
		xval[i] = x[i]
	}

}
